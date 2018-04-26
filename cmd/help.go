package cmd

import "fmt"

func ShowHelp()  {
	fmt.Println("Welcome use Ttpush!")
	fmt.Println("===================")
	fmt.Println("")
	fmt.Println("command:")
	fmt.Printf("\t%s\n","-c :")
	fmt.Printf("\t\t%s\n","input config file path,content like this:")
	fmt.Printf("\t\t%s\n",`{
		  "server_port":"8080",
		  "max_ttl":2419200,
		  "api_key":"AAAA_1dLSps:APA91......ZHrCUioe-vx6wFvDXfnoh9h",
		  "notify_callback":"http://localhost:8000/fcm/notify",
		  "log_file":"/tmp/",
		  "proxy":"",
		  "notification":{
			"title":"",
			"body":"",
			"icon":"http://",
			"uri":"http://"
		  }
		}`)
	fmt.Println("")
	fmt.Printf("\t%s\n","-h :")
	fmt.Printf("\t\t%s\n","help command and list commands")
	fmt.Println("")
	fmt.Printf("\t%s\n","-d :")
	fmt.Printf("\t\t%s\n","true or false set It's daemon?")
}