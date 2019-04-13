<?php

function helloWorld() {
    $thing = 2;

    if ($thing > 10) {
        print('Hello');
    } else if ($thing > 5) {
        print('World');
    } else if ($thing > 2) {
        print('Yo');
    }
}
