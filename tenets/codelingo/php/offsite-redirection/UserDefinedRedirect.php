<?php

/**
 * @Route("/login")
 */
class Controller
{
    public function viewAction(Request $request)
    {
        $url = $request->query->get("url", '');
        $logoutUrl = "/logout";
        if ($url === $logoutUrl || $url === '') {
            $url = "/";
        }

        $user = $this->getUserFromSession();
        if ($user !== null) {
            // they are logged in so redirect them to a logged in page
            return $this->redirect('' . offsiteCheck($url));
        } else if ($user === "arbitrary") {
            return $this->redirect("");
        } else {
            return $this->redirect($url); // ISSUE
        }
    }

    public function getUserFromSession() {
        return "currentUser";
    }
    
    public function redirect($url) {
        echo "redirecting to " . $url;
    }
}

// stand in for Symfony request
class Request
{
    public function __construct() {
        $this->query = new Query();
    }
}

Class Query {
    public function get () {
        return "arbitraryURL";
    }
}

function ExecuteSQL ($query) {
	return "SQL result\n";
}

function offsiteCheck($url) {
	if ($url) {
        $urlParts = parse_url($url);
    }
    if ($url && !preg_match('/logout/', $url) && !isset($urlParts['host']) && substr(str_replace('\\', '', $url), 0, 2) !== '//' && substr($url, 0, 2) !== '\\\\') {
        // the url wasn't an offsite redirect
        echo $url . "\n";
    } else {
        // suspect url, forward to the index page
        echo '/index' . "\n";
    }
}

$req = new Request();

$ctrl = new Controller();
$ctrl->viewAction($req);



offsiteCheck("https://www.badsite.com");
offsiteCheck("https://www.goodsite.com");
offsiteCheck("/localpath");
offsiteCheck("data:text/html;charset=UTF-8,<html><script>window.location=\"https://www.badsite.com\"</script></html>");