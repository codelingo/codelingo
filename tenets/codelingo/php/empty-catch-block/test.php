<?php
// From http://php.net/manual/en/language.exceptions.php
function inverse($x) {
    if (!$x) {
        throw new Exception('Division by zero.');
    }
    return 1/$x;
}

try {
	echo "In silly try block\n";
    echo inverse(0) . "\n";
} catch (Exception $e) {
}

try {
    echo inverse(5) . "\n";
    echo inverse(0) . "\n";
} catch (Exception $e) {
    echo 'Caught exception: ',  $e->getMessage(), "\n";
}

// Continue execution
echo "Hello World\n";
?>