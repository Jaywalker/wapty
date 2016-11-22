package webroot

import (
	"log"
)

func init() {
	log.Println("Got here")
	webFiles[""] = indexHtml
}

const indexHtml = `<HTML>
	<HEAD>
	</HEAD>
	<BODY>
		<button id="forwardOriginal" type="button" onclick="clickhandler();">Forward Original</button>
		<button id="forwardModified" type="button" onclick="clickhandler();">Forward Modified</button>
		<button id="drop" type="button" onclick="clickhandler();">Drop</button>
		<button id="provideResponse" type="button" onclick="clickhandler();">Provide Response</button>
		<label><input type="checkbox" id="interceptToggle" value="Intercept" onclick="clickhandler();">Intercept</label>

		<button id="menu" type="button">Action</button>
		<br>
		<br>

		<textarea id="proxybuffer" name="Text1" cols="100" rows="40"></textarea>

	</BODY>
	<SCRIPT> 
var intercept = {
	EDITORCHANNEL: "proxy/intercept/editor",
	SETTINGSCHANNEL: "proxy/intercept/options",
	HISTORYCHANNEL: "proxy/httpHistory",
}
var waptyServer= new WebSocket("ws://localhost:8081/ws");
waptyServer.onopen = function(event){
	console.log("WebSocket connected");
	var msg = {
		Action: "intercept",
		Channel: intercept.SETTINGSCHANNEL
	}
	waptyServer.send(JSON.stringify(msg));

}
waptyServer.onmessage = function(event){
	//	console.log(event.data);
	msg = JSON.parse(event.data);
	console.log(msg);
	switch (msg.Channel){
		case intercept.EDITORCHANNEL:
			//if ('Payload' in msg){
			console.log(atob(msg.Payload));
			document.getElementById("proxybuffer").value=atob(msg.Payload);
			//}
			break;
		case intercept.SETTINGSCHANNEL:
			switch (msg.Action){
				case "intercept":
					document.getElementById("interceptToggle").checked = msg.Args[0] === "true";
			}
			break;
		case intercept.HISTORYCHANNEL:
			switch (msg.Action){
				case "metaData":
					var metaData = JSON.parse(msg.Args[0])
					console.log("Metadata for request " + metaData.Id + " received:");
					console.log(metaData)
			}
			break;
	}
}

function clickhandler(){
	switch (event.target.id){
		case "forwardOriginal":
			var msg = {
				Action: "forward",
				Channel: intercept.EDITORCHANNEL
			}
			document.getElementById("proxybuffer").value="";
			waptyServer.send(JSON.stringify(msg));
			break;
		case "forwardModified":
			var payload = btoa(document.getElementById("proxybuffer").value);
			var msg = {
				Action: "edit",
				Channel: intercept.EDITORCHANNEL,
				Payload: payload
			}
			document.getElementById("proxybuffer").value="";
			waptyServer.send(JSON.stringify(msg));
			break;

			break;
		case "drop":
			break;
		case "provideResponse":
			break;
		case "interceptToggle":
			var msg = {
				Action: "intercept",
				Channel: intercept.SETTINGSCHANNEL,
				Args: [""+document.getElementById("interceptToggle").checked]
			}
			waptyServer.send(JSON.stringify(msg));
			break;
		default:
			console.log("unknown event")
	}
}
	</SCRIPT>
</HTML>
`