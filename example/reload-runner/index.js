const WebSocket = require('ws');
const schema = require('enigma.js/schemas/12.20.0.json');
var enigma = require('enigma.js')

// create a new session:
const session = enigma.create({
	schema,
	url: 'ws://localhost:19076/app/engineData',
	createSocket: url => new WebSocket(url),
});

// bind traffic events to log what is sent and received on the socket:
session.on('traffic:sent', data => console.log('sent:', data));
session.on('traffic:received', data => console.log('received:', data));

// open the socket and eventually receive the QIX global API, and then close
// the session:


var global
var app;
var reloadRequestId;
var appId = "app6.qvf";
session.open()

	.then((_global) => {
		global = _global;
		console.log('Creating session app');
		return global.createApp(appId).then( (appInfo) => {
			return global.openDoc(appInfo.qAppId)
		}).catch( (err) => {
			return global.openDoc(appId)
		});
	})
	.then((_app) => {
		console.log('Creating connection');
		app = _app;
		return app.createConnection({
			qType: 'Postgres',
			qName: 'postgresgrpc2',
			qConnectionString: 'CUSTOM CONNECT TO "provider=Postgres;host=postgres-database;port=5432;database=postgres;user=postgres;password=postgres"',
			//qConnectionString: 'CUSTOM CONNECT TO "provider=Postgres;host=selun-gwe.qliktech.com;port=5433;database=postgres;user=postgres;password=postgres"',
			qUserName: 'postgres',
			qPassword: 'postgres'
		})
	})
	.then(() => {
		console.log('Setting script');
		const script = `
			lib connect to 'postgresgrpc2';		
			Airports:						
			sql select * from airports;
		`;
		return app.setScript(script);
	})
	.then(() => {
		console.log('Reloading');
		var reloadPromise = app.doReload();
		reloadRequestId = reloadPromise.requestId;
		return reloadPromise;
	})
	.then(() => {
		return global.getProgress(reloadRequestId)
	})
	.then((progress) => {

	})
	.then(() => {
		console.log('Fetching Table sample');
		return app.getTableData(-1, 50, true, "Airports")
	})
	.then((tableData) => {
		//Convert table grid into a string using some functional magic
		var tableDataAsString = tableData.map((row) => row.qValue.map((value) => value.qText).reduce((left, right) => left + "\t" + right)).reduce((row1, row2) => row1 + "\n" + row2)
		console.log(tableDataAsString);
	})
	.then(() => session.close())
	.then(() => console.log('Session closed'))
	.catch(err => console.log('Something went wrong :(', err))
