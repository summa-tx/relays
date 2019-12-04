package keeper

import (
	"encoding/json"
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

type LinkTest struct{}

type ChainTest struct{}

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

type KeeperTestCases struct {
	LinkTestCases   []LinkTest      `json:"link"`
	HeaderTestCases HeaderTestCases `json:"header"`
	ChainTestCases  []ChainTest     `json:"chain"`
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

func (suite *KeeperSuite) InitTestContext(mainnet, isCheckTx bool) {
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

	suite.Context = ctx
	suite.Keeper = keeper
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

	keeperSuite.InitTestContext(true, false)

	suite.Run(t, keeperSuite)
}

func (s *KeeperSuite) TestGetPrefixStore() {
	prefStore := s.Keeper.getPrefixStore(s.Context, "toast-")
	store := s.Context.KVStore(s.Keeper.storeKey)

	expected := []byte{0xff}

	prefStore.Set([]byte("1"), expected)
	actual := store.Get([]byte("toast-1"))

	s.Equal(expected, actual)
}
