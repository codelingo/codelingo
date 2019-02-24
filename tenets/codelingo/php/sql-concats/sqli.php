<?php

class Test{
    public function A($param) {
        $query  = "INSERT INTO some_table (some_value) VALUES ('".$param."')";
        $this->dbexec($query);
    }
    public function B($param){
        $param = $this->sanitise($param);
        $query  = "INSERT INTO some_table (some_value) VALUES ('".$param."')";
        $this->dbexec($query);
    }
    public function C($param){
        $param = $this->D($param);
        $query  = "INSERT INTO some_table (some_value) VALUES ('".$param."')";
        $this->dbexec($query);
    }
    public function D($param){
        $param = $this->sanitise($param);
        return $param;
    }
    public function E($param){
        $param = $this->F($param);
        $query  = "INSERT INTO some_table (some_value) VALUES ('".$param."')";
        $this->dbexec($query);
    }
    public function F($param){
        return $param;
    }
    public function dbexec($query) {
        echo $query."\n";
    }
    public function sanitise($unsafe_value){
        $safe_value = substr($unsafe_value,2);
        return $safe_value;
    }
    public function httpGetRequest(){
        $userinput_unsafe = "unsafe value";
        return $userinput_unsafe;
    }
}

$test = new Test();
$userinput = $test->httpGetRequest();
$safeinput = $test->sanitise($userinput);
$test->A($safeinput);
$test->A($userinput);
$test->B($userinput);
$test->C($userinput);
$test->E($userinput);

