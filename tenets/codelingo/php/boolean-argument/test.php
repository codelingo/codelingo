<?php

function takes_bool(bool $input) // ISSUE
{
    if ($input) {
        echo "true";
    } else {
        echo "false";
    }
}

function does_not_take_bool(int $input)
{
    if ($input == 10) {
        echo "true";
    } else {
        echo "false";
    }
}

takes_bool(true);
takes_bool(false);

does_not_take_bool(10);
?>