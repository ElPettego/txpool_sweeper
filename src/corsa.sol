// SPDX-License-Identifier: MIT 

pragma solidity ^0.8.19;

import "./contracts/IPancakeRouter01.sol";
import "./contracts/IPancakeRouter02.sol";
import "./contracts/IERC20.sol";

contract Correre {
    address internal constant PANCAKE_ROUTER_V2_ADDRESS = 0x10ED43C718714eb63d5aA57B78B54704E256024E;

    IPancakeRouter02 public pcsRouter;
    uint constant MAX_UINT = 2**256 - 1;
    address payable owner;

    event Received(address sender, uint amount);

    constructor() {
        pcsRouter = IPancakeRouter02(PANCAKE_ROUTER_V2_ADDRESS);
        owner = payable(msg.sender);

    }

    modifier onlyOwner {
        require(
            msg.sender == owner, "__u_cant_do_this_bruv__"
        );
        _;
    }

    receive() external payable {
        emit Received(msg.sender, msg.value);
    }

    function buyToken(uint ethAmount, address tokenAddress) public payable onlyOwner {
        require(ethAmount <= address(this).balance, "__too_lil_eth__");
        IERC20 token = IERC20(tokenAddress);
        if (token.allowance(address(this), PANCAKE_ROUTER_V2_ADDRESS) < 1) {
            require(token.approve(PANCAKE_ROUTER_V2_ADDRESS, MAX_UINT), "__approve_failed__");
        }
        address[] memory path = new address[](2);
        path[0] = pcsRouter.WETH();
        path[1] = tokenAddress;
        pcsRouter.swapExactETHForTokens{value: ethAmount}(0, path, address(this), block.timestamp + 60);
    }

    function sellToken(address tokenAddress) public payable onlyOwner {
        IERC20 token = IERC20(tokenAddress);
        address[] memory path = new address[](2);
        path[0] = tokenAddress;
        path[1] = pcsRouter.WETH();
        uint tokenBalance = token.balanceOf(address(this));
        pcsRouter.swapExactTokensForETH(tokenBalance, 0, path, address(this), block.timestamp + 60);
    } 

    function withdraw() public payable onlyOwner {
        owner.transfer(address(this).balance);
    }

    function doThaShit(uint ethAmount, address tokenAddress) public payable onlyOwner {
        buyToken(ethAmount, tokenAddress);
        sellToken(tokenAddress);
    }

    function withdrawToken(address tokenAddress) public payable onlyOwner {
        IERC20 token = IERC20(tokenAddress);
        uint tokenBalance = token.balanceOf(address(this));
        require(tokenBalance > 0, "__no_token__");
        token.transfer(owner, tokenBalance);
    }
}
