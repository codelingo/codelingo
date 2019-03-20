<?php
namespace test;
class Foo {
    public function aMethod(bool $param) {
        echo $param . "\n";
        $this->aMethod(TRUE);
        $this->test();
    }

    public function test() {
        $this->aMethod(TRUE);
        $this->called();
    }

    private function called() {
        echo "called";
    }

    private function not_called1() {
        echo "not called";
    }

    private function not_called2() {
        echo "not called";
    }
}