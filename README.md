ddnssrv
=======

ddnssrv is a simple dynamic DNS listener that accepts DNS updates in the same format as the DynDNS.org servers.

The server was developed to be used in combination with [https://github.com/marpie/cfup](cfup.py).

Sample Configuration
--------------------

	{
	  "Users": [
	    {
	      "Username": "Default-User",
	      "Password": "MyPassword",
	      "Domains": [
	        "www.example.com"
	      ]
	    }
	  ]
	}
