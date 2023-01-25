// SPDX-License-Identifier: GPL-2.0-or-later
pragma solidity =0.7.6;
pragma abicoder v2;

// import "@uniswap/v3-periphery/contracts/libraries/TransferHelper.sol";
import "@uniswap/v3-periphery/contracts/interfaces/ISwapRouter.sol";

interface IERC20 {
    function balanceOf(address account) external view returns (uint256);

    function transfer(address recipient, uint256 amount)
        external
        returns (bool);

    function approve(address spender, uint256 amount) external returns (bool);
}

contract SingleSwap {
    address public constant routerAddress =
        0xE592427A0AEce92De3Edee1F18E0157C05861564;
    ISwapRouter public immutable swapRouter = ISwapRouter(routerAddress);

    address public constant WETH = 0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6;
    address public constant USDC = 0x07865c6E87B9F70255377e024ace6630C1Eaa37F;

    IERC20 public wETHToken = IERC20(WETH);

    // For this example, we will set the pool fee to 0.3%.
    uint24 public constant poolFee = 3000;

    constructor() {}

    function swapExactInputSingle(uint256 amountIn)
        external
        returns (uint256 amountOut)
    {
        /* 
        /  metamask or some outside wallet must first approve the funds for the contract to spend
        /  should be executed by arbiter 
        */
        // TransferHelper.safeTransferFrom(
        //     WETH,
        //     msg.sender,
        //     address(this),
        //     amountIn
        // );

        // TransferHelper.safeApprove(WETH, address(swapRouter), amountIn);
        // replace above with below
        wETHToken.approve(address(swapRouter), amountIn);

        ISwapRouter.ExactInputSingleParams memory params = ISwapRouter
            .ExactInputSingleParams({
                tokenIn: WETH,
                tokenOut: USDC,
                fee: poolFee,
                recipient: address(this), // changed from msg.sender
                deadline: block.timestamp,
                amountIn: amountIn,
                amountOutMinimum: 0,
                sqrtPriceLimitX96: 0
            });

        // The call to `exactInputSingle` executes the swap.
        amountOut = swapRouter.exactInputSingle(params);
    }

    // probably not gonna use this
    function swapExactOutputSingle(uint256 amountOut, uint256 amountInMaximum)
        external
        returns (uint256 amountIn)
    {
        // same with above. must be approved to spend tokens by outside wallet
        // TransferHelper.safeTransferFrom(
        //     WETH,
        //     msg.sender,
        //     address(this),
        //     amountInMaximum
        // );

        // TransferHelper.safeApprove(WETH, address(swapRouter), amountInMaximum);
        // replace above with below
        wETHToken.approve(address(swapRouter), amountInMaximum);

        ISwapRouter.ExactOutputSingleParams memory params = ISwapRouter
            .ExactOutputSingleParams({
                tokenIn: WETH,
                tokenOut: USDC,
                fee: poolFee,
                recipient: address(this), // changed from msg.sender
                deadline: block.timestamp,
                amountOut: amountOut,
                amountInMaximum: amountInMaximum,
                sqrtPriceLimitX96: 0
            });

        // Executes the swap returning the amountIn needed to spend to receive the desired amountOut.
        amountIn = swapRouter.exactOutputSingle(params);

        // For exact output swaps, the amountInMaximum may not have all been spent.
        // If the actual amount spent (amountIn) is less than the specified maximum amount, we must refund the msg.sender and approve the swapRouter to spend 0.
        if (amountIn < amountInMaximum) {
            // TransferHelper.safeApprove(WETH, address(swapRouter), 0);
            // TransferHelper.safeTransfer(
            //     WETH,
            //     address(this), // change from msg.sender. switch back once authorizing from external wallet
            //     amountInMaximum - amountIn
            // );

            // replace TransferHelper with below
            wETHToken.approve(address(swapRouter), 0);
            wETHToken.transfer(address(this), amountInMaximum);
        }
    }
}
