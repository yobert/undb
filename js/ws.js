var tools = require('github.com/yobert/pkunk/js/tools');

function UndbSocket(store, uri) {
	var ws = new WebSocket(uri);
	var ws_id = 'server-ws';

	this.sock = ws;

	var firstmessage = false;
	var t = this;

	var softerror;

	var onchange = function(op) {
		if(op.changesource && op.changesource == ws_id)
			return;

		var o = {
			Method: op.Method,
			Path: op.Path,
			Id: op.Id,
			Type: op.Type,
			Values: op.Values
		};
		var msg = JSON.stringify([o]);
		//console.log('send: '+msg);
		ws.send(msg);
	};

	ws.onopen = function() {
		//console.log("websocket open");
		store.addListener('CHANGE', onchange);
	};

	ws.onclose = function() {
		store.removeListener('CHANGE', onchange);
		//console.log("websocket closed");
		if(t.onClose)
			t.onClose(softerror);

		softerror = undefined;
	};

	//ws.onerror = function(err) {
	//	console.log("websocket error");
	//};

	ws.onmessage = function(msg) {
		//console.log('recv: '+msg.data)

		var oplist = JSON.parse(msg.data);
		if(oplist && oplist.Error) {
			softerror = oplist.Error;
			return;
		}

		var err;
		for(var i in oplist) {
			err = store.Exec(oplist[i], ws_id);

			if(err)
				throw(err);
		}

		if(!firstmessage) {
			firstmessage = true;

			if(t.onFirstMessage)
				t.onFirstMessage();
		}
	};
}

module.exports = {
	UndbSocket: UndbSocket
};
