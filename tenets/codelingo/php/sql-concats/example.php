<?php
function ExecuteSQL ($query) {
	return "SQL result\n";
}
function Main() {
	$userTable = "USERS";
	$userNamesSQL = "SELECT name FROM " . $userTable; // ISSUE
	$userNames = ExecuteSQL($userNamesSQL);
	print $userNames;

	$sql = "SELECT cost FROM ";
	$sql .= "PRODUCTS"; // ISSUE
	$products = ExecuteSQL($sql);
	print $$products;
}
Main();
php>