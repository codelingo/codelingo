<?php
class Foo {
    public static function aStaticMethod(bool $param) {
        echo $param . "\n";
        return $param;
    }
    public static function anotherStaticMethod(bool $param) {
        echo $param . "\n";
    }

}

$a = Foo::aStaticMethod(TRUE);
$b = Foo::anotherStaticMethod(TRUE);