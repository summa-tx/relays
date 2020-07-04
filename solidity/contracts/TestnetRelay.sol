pragma solidity ^0.5.10;

/** @title TestnetRelay */
/** @author Summa (https://summa.one) */

import {OnDemandSPV} from "./OnDemandSPV.sol";
import {TypedMemView} from "@summa-tx/bitcoin-spv-sol/contracts/TypedMemView.sol";

contract TestnetRelay is OnDemandSPV {

    constructor(
        bytes memory _genesisHeader,
        uint256 _height,
        bytes32 _periodStart,
        uint256 _firstID
    ) OnDemandSPV(
        _genesisHeader,
        _height,
        _periodStart,
        _firstID
    ) public {return ;}

    function _addHeadersWithRetarget(
        bytes memory, // _oldPeriodStartHeader,
        bytes memory _oldPeriodEndHeader,
        bytes memory _headers
    ) internal returns (bool) {
        bytes29 _oldEnd = _oldPeriodEndHeader.ref(0).tryAsHeader();
        bytes29 _headersView = _headers.ref(0).tryAsHeaderArray();

        require(
            _oldEnd.notNull() && _headersView.notNull(),
            "Bad args. Check header and array byte lengths."
        );
        return _addHeaders(_oldEnd, _headersView, true);
    }

    /// @notice             Adds headers to storage after validating
    /// @dev                We check integrity and consistency of the header chain
    /// @param  _anchor     The header immediately preceeding the new chain
    /// @param  _headers    A tightly-packed list of new 80-byte Bitcoin headers to record
    /// @return             True if successfully written, error otherwise
    function _addHeaders(bytes29 _anchor, bytes29 _headers, bool _internal) internal returns (bool) {
        /// Extract basic info
        bytes32 _previousDigest = _anchor.hash256();
        uint256 _anchorHeight = _findHeight(_previousDigest);  /* NB: errors if unknown */
        uint256 _target = _headers.indexHeaderArray(0).target();

        require(
            _internal || _anchor.target() == _target,
            "Unexpected retarget on external call"
        );

        /*
        NB:
        1. check that the header has sufficient work
        2. check that headers are in a coherent chain (no retargets, hash links good)
        3. Store the block connection
        4. Store the height
        */
        uint256 _height;
        bytes32 _currentDigest;
        for (uint256 i = 0; i < _headers.len() / 80; i += 1) {
            bytes29 _header = _headers.indexHeaderArray(i);
            _height = _anchorHeight.add(i + 1);
            _currentDigest = _header.hash256();

            /*
            NB:
            if the block is already authenticated, we don't need to a work check
            Or write anything to state. This saves gas
            */
            if (previousBlock[_currentDigest] == bytes32(0)) {
                require(
                    TypedMemView.reverseUint256(uint256(_currentDigest)) <= _target,
                    "Header work is insufficient"
                );
                previousBlock[_currentDigest] = _previousDigest;
                if (_height % HEIGHT_INTERVAL == 0) {
                    /*
                    NB: We store the height only every 4th header to save gas
                    */
                    blockHeight[_currentDigest] = _height;
                }
            }

            /* NB: we do still need to make chain level checks tho */
            require(_header.target() == _target, "Target changed unexpectedly");
            require(_header.checkParent(_previousDigest), "Headers do not form a consistent chain");

            _previousDigest = _currentDigest;
        }

        emit Extension(
            _anchor.hash256(),
            _currentDigest);
        return true;
    }
}
