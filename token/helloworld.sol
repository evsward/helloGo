pragma solidity ^0.4.0;

contract helloworld {
    string content;

    function helloworld(string _str) public {
        content = _str;
    }

    function getContent() constant public returns (string){
        return content;
    }
}
