pragma solidity ^0.4.16;

contract enumsTest {
    enum ColorEnums {Red, White, Black}
    ColorEnums color;
    ColorEnums constant defaultColor = ColorEnums.White;

    function setBlack() public {
        this.color = ColorEnums.Black;
    }

    function getColor() view public returns (ColorEnums) {
        return color;
    }

    function getDefault() public pure returns (uint) {
        return uint(defaultColor);
    }
}
