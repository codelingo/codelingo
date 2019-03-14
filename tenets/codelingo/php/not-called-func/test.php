<?php
class Foo {
    public function aMethod(bool $param) {
        echo $param . "\n";
        $this->test($param);
    }

    public function test($param) {
        echo $param . "\n";
    }


    public function not_called() {
        echo "not called";
    }
}