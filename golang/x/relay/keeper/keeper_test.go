package keeper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/stretchr/testify/suite"

	"github.com/summa-tx/relays/golang/x/relay/types"
)

type NamedCase interface {
	Name() string
}

type Case struct {
	NamedCase
	Comment string `json:"comment"`
}

/***** LINK TEST CASES *****/
type IsAncestorTestCase struct {
	Digest   types.Hash256Digest `json:"digest"`
	Ancestor types.Hash256Digest `json:"ancestor"`
	Limit    uint32              `json:"limit"`
	Output   bool                `json:"output"`
}

type IsAncestor struct {
	TestCases []IsAncestorTestCase `json:"testCases"`
}

type FindAncestorTestCase struct {
	Digest types.Hash256Digest `json:"digest"`
	Offset uint32              `json:"offset"`
	Error  int                 `json:"error"`
	Output types.Hash256Digest `json:"output"`
}

type FindAncestor struct {
	TestCases []FindAncestorTestCase `json:"testCases"`
}

type LinkTestCases struct {
	IsAncestor   IsAncestor   `json:"isAncestor"`
	FindAncestor FindAncestor `json:"findAncestor"`
}

/***** CHAIN TEST CASES *****/
type MostRecentCATestCase struct {
	Ancestor types.Hash256Digest `json:"ancestor"`
	Left     types.Hash256Digest `json:"left"`
	Right    types.Hash256Digest `json:"right"`
	Limit    uint32              `json:"limit"`
	Output   bool                `json:"output"`
}

type IsMostRecentCA struct {
	Orphan            types.BitcoinHeader    `json:"orphan"`
	OldPeriodStart    types.BitcoinHeader    `json:"oldPeriodStart"`
	Genesis           types.BitcoinHeader    `json:"genesis"`
	PreRetargetChain  []types.BitcoinHeader  `json:"preRetargetChain"`
	PostRetargetChain []types.BitcoinHeader  `json:"postRetargetChain"`
	TestCases         []MostRecentCATestCase `json:"testCases"`
}

type HeaviestTestCase struct {
	Ancestor    types.Hash256Digest `json:"ancestor"`
	CurrentBest types.Hash256Digest `json:"currentBest"`
	NewBest     types.Hash256Digest `json:"newBest"`
	Limit       uint32              `json:"limit"`
	Error       int                 `json:"error"`
	Output      types.Hash256Digest `json:"output"`
}

type HeaviestFromAncestor struct {
	Orphan    types.BitcoinHeader   `json:"orphan"`
	BadHeader types.BitcoinHeader   `json:"badHeader"`
	Genesis   types.BitcoinHeader   `json:"genesis"`
	Headers   []types.BitcoinHeader `json:"headers"`
	TestCases []HeaviestTestCase    `json:"testCases"`
}

type NewHeaviestTestCase struct {
	BestKnownDigest types.Hash256Digest `json:"bestKnownDigest"`
	Ancestor        types.Hash256Digest `json:"ancestor"`
	CurrentBest     types.RawHeader     `json:"currentBest"`
	NewBest         types.RawHeader     `json:"newBest"`
	Limit           uint32              `json:"limit"`
	Error           int                 `json:"error"`
	Output          string              `json:"output"`
}

type ChainTestCases struct {
	IsMostRecentCA       IsMostRecentCA        `json:"isMostRecentCommonAncestor"`
	HeaviestFromAncestor HeaviestFromAncestor  `json:"heaviestFromAncestor"`
	MarkNewHeaviest      []NewHeaviestTestCase `json:"markNewHeaviest"`
}

/***** HEADER TEST CASES *****/
type IngestCase struct {
	Case
	Headers   []types.BitcoinHeader `json:"headers"`
	Anchor    types.BitcoinHeader   `json:"anchor"`
	Internal  bool                  `json:"internal"`
	IsMainnet bool                  `json:"isMainnet"`
	Output    sdk.CodeType          `json:"output"`
}

type DiffChangeCase struct {
	Case
	Headers        []types.BitcoinHeader `json:"headers"`
	PrevEpochStart types.BitcoinHeader   `json:"prevEpochStart"`
	Anchor         types.BitcoinHeader   `json:"anchor"`
	Output         sdk.CodeType          `json:"output"`
}

type CompareCase struct {
	Case
	Full      sdk.Uint `json:"full"`
	Truncated sdk.Uint `json:"truncated"`
	Output    bool     `json:"output"`
}

type HeaderTestCases struct {
	ValidateDiffChange []DiffChangeCase `json:"validateDifficultyChange"`
	ValidateChain      []IngestCase     `json:"validateHeaderChain"`
	CompareTargets     []CompareCase    `json:"compareTargets"`
}

/***** Validator TEST CASES *****/
type ValidateProofTestCase struct {
	Proof     types.SPVProof      `json:"proof"`
	BestKnown types.Hash256Digest `json:"bestKnown"`
	LCA       types.Hash256Digest `json:"lca"`
	Error     int                 `json:"error"`
	Output    bool                `json:"output"`
}

type ValidatorTestCases struct {
	ValidateProof []ValidateProofTestCase `json:"validateProof"`
}

/***** Request TEST CASES *****/
type CheckRequestTestCase struct {
	InputIdx  uint8          `json:"inputIndex"`
	OutputIdx uint8          `json:"outputIndex"`
	Vin       types.HexBytes `json:"vin"`
	Vout      types.HexBytes `json:"vout"`
	RequestID types.HexBytes `json:"requestID"`
	Error     int            `json:"error"`
	Output    bool           `json:"output"`
}

type RequestTestCases struct {
	EmptyRequest  types.ProofRequest     `json:"emptyRequest"`
	CheckRequests []CheckRequestTestCase `json:"checkRequests"`
}

/***** KEEPER TEST CASES *****/
type KeeperTestCases struct {
	LinkTestCases      LinkTestCases      `json:"link"`
	HeaderTestCases    HeaderTestCases    `json:"header"`
	ChainTestCases     ChainTestCases     `json:"chain"`
	ValidatorTestCases ValidatorTestCases `json:"validator"`
	RequestTestCases   RequestTestCases   `json:"requests"`
}

type KeeperSuite struct {
	suite.Suite
	Fixtures KeeperTestCases
	Context  sdk.Context
	Keeper   Keeper
}

func (c Case) Name() string {
	return c.Comment
}

func logIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func logIfTestCaseError(tc NamedCase, err sdk.Error) {
	if err != nil {
		log.Printf("Unexpected Error\nIn case: %s\n%s\n", tc.Name(), err.Error())
	}
}

func (s *KeeperSuite) InitTestContext(mainnet, isCheckTx bool) {
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

	relayKey := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(relayKey, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err.Error())
	}

	cdc := codec.New()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "relayTestChain"}, isCheckTx, tmlog.NewNopLogger())
	keeper := NewKeeper(relayKey, cdc, mainnet)

	s.Context = ctx
	s.Keeper = keeper
}

func (s *KeeperSuite) SetupTest() {
	s.InitTestContext(true, false)
}

// Runs the whole test suite
func TestKeeper(t *testing.T) {
	jsonFile, err := os.Open("../../../../testVectors.json")
	logIfError(err)
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	logIfError(err)

	var fixtures KeeperTestCases
	err = json.Unmarshal([]byte(byteValue), &fixtures)
	logIfError(err)

	keeperSuite := new(KeeperSuite)
	keeperSuite.Fixtures = fixtures

	suite.Run(t, keeperSuite)
}

func (suite *KeeperSuite) SDKNil(e sdk.Error) {
	var msg string
	if e != nil {
		msg = e.Error()
	}
	suite.Nil(e, msg)
}

func (suite *KeeperSuite) EqualError(e sdk.Error, code int) {
	var msg string
	if e.Code() != sdk.CodeType(code) {
		msg = fmt.Sprintf("%sExpected: %d\n", e.Error(), code)
	}
	suite.Equal(e.Code(), sdk.CodeType(code), msg)
}

func (s *KeeperSuite) TestGetPrefixStore() {
	prefStore := s.Keeper.getPrefixStore(s.Context, "toast-")
	store := s.Context.KVStore(s.Keeper.storeKey)

	expected := []byte{0xff}

	prefStore.Set([]byte("1"), expected)
	actual := store.Get([]byte("toast-1"))

	s.Equal(expected, actual)
}

func (s *KeeperSuite) TestSetGenesisState() {
	genesis := s.Fixtures.HeaderTestCases.ValidateDiffChange[0].Anchor
	epochStart := s.Fixtures.HeaderTestCases.ValidateDiffChange[0].PrevEpochStart
	err := s.Keeper.SetGenesisState(s.Context, genesis, epochStart)
	s.SDKNil(err)

	gen, err := s.Keeper.GetRelayGenesis(s.Context)
	s.SDKNil(err)
	s.Equal(genesis.HashLE, gen)

	lca, err := s.Keeper.GetLastReorgLCA(s.Context)
	s.SDKNil(err)
	s.Equal(genesis.HashLE, lca)

	best, err := s.Keeper.GetBestKnownDigest(s.Context)
	s.SDKNil(err)
	s.Equal(genesis.HashLE, best)

	err = s.Keeper.SetGenesisState(s.Context, genesis, epochStart)
	s.Equal(types.AlreadyInit, err.Code())
}
