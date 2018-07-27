<?php 
for ($x = 0; $x <= 10; $x++) {
    echo "$x\n";
} 

$i = 0;
do {
    echo $i;
} while ($i > 0);

$j = 1;
while ($j <= 10) {
    echo $j++;
}


for ($x = 0; $x <= 10; $x++) {
    echo sizeof($x);
    echo "$x\n";
} 

$i = 0;
do {
    echo count($i);
    echo $i;
} while ($i > 0);

$j = 1;
while ($j <= 10) {
    echo sizeof($j);
    echo $j++;
}
?>