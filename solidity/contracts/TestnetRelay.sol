pragma solidity ^0.5.10;

/** @title TestnetRelay */
/** @author Summa (https://summa.one) */

import {OnDemandSPV} from "./OnDemandSPV.sol";

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
        return _addHeaders(_oldPeriodEndHeader, _headers);
    }

    function _addHeaders(bytes memory _anchor, bytes memory _headers, bool) internal returns (bool) {
        return _addHeaders(_anchor, _headers);
    }

    /// @notice             Adds headers to storage after validating
    /// @dev                We check integrity and consistency of the header chain
    /// @param  _anchor     The header immediately preceeding the new chain
    /// @param  _headers    A tightly-packed list of new 80-byte Bitcoin headers to record
    /// @return             True if successfully written, error otherwise
    function _addHeaders(bytes memory _anchor, bytes memory _headers) internal returns (bool) {
        uint256 _height;
        bytes memory _header;
        bytes32 _currentDigest;
        bytes32 _previousDigest = _anchor.hash256();

        /* uint256 _target = _headers.slice(0, 80).extractTarget(); */
        uint256 _anchorHeight = _findHeight(_previousDigest);  /* NB: errors if unknown */

        /* require(
            _internal || _anchor.extractTarget() == _target,
            "Unexpected retarget on external call"); */
        require(_headers.length % 80 == 0, "Header array length must be divisible by 80");

        /*
        NB:
        1. check that the header has sufficient work
        2. check that headers are in a coherent chain (no retargets, hash links good)
        3. Store the block connection
        4. Store the height
        */
        for (uint256 i = 0; i < _headers.length / 80; i = i.add(1)) {
            _header = _headers.slice(i.mul(80), 80);
            _height = _anchorHeight.add(i + 1);
            _currentDigest = _header.hash256();

            /*
            NB:
            if the block is already authenticated, we don't need to a work check
            Or write anything to state. This saves gas
            */
            if (previousBlock[_currentDigest] == bytes32(0)) {
                require(
                    abi.encodePacked(_currentDigest).reverseEndianness().bytesToUint() <= _header.extractTarget(),
                    "Header work is insufficient");
                previousBlock[_currentDigest] = _previousDigest;
                if (_height % HEIGHT_INTERVAL == 0) {
                    /*
                    NB: We store the height only every 4th header to save gas
                    */
                    blockHeight[_currentDigest] = _height;
                }
            }

            /* NB: we do still need to make chain level checks tho */
            /* require(_header.extractTarget() == _target, "Target changed unexpectedly"); */
            require(_header.validateHeaderPrevHash(_previousDigest), "Headers do not form a consistent chain");

            _previousDigest = _currentDigest;
        }

        emit Extension(
            _anchor.hash256(),
            _currentDigest);
        return true;
    }
}
