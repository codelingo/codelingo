<?php
if ($foo = 'bar') { // ISSUE: possible typo
    echo "foob";
}
if ($baz = 0) { // ISSUE: always false
    echo "baz";
}
if (true) {
    echo "true";
}
$thing = "ting";
?>