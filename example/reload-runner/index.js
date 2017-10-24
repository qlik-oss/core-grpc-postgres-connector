const WebSocket = require('ws');
const schema = require('enigma.js/schemas/12.20.0.json');
var enigma = require('enigma.js')

// create a new session:
const session = enigma.create({
	schema,
	url: 'ws://localhost:9076/app/engineData',
	createSocket: url => new WebSocket(url),
});

var global
var app;
var reloadRequestId;
var appId = "reloadapp.qvf";
var connectionId;
var trafficLog = false;

if (trafficLog) {
	// bind traffic events to log what is sent and received on the socket:
 	session.on('traffic:sent', data => console.log('sent:', data));
	session.on('traffic:received', data => console.log('received:', data));
}

session.open()
	.then((_global) => {
		global = _global;
		console.log('Creating/opening app');
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
			qType: 'postgres-grpc-connector', //the name we defined as a parameter to engine in our docker-compose.yml
			qName: 'postgresgrpc',
			qConnectionString: 'CUSTOM CONNECT TO "provider=postgres-grpc-connector;host=postgres-database;port=5432;database=postgres;user=postgres;password=postgres"', //the connection string inclues both the provide to use and parameters to it.
			qUserName: 'postgres',
			qPassword: 'postgres'
		})
	})
	.then((_connectionId) => {
		connectionId = _connectionId;
		console.log('Setting script');
		const script = `
			lib connect to 'postgresgrpc';		
			Airports:						
			sql select rowID,Airport,City,Country,IATACode,ICAOCode,Latitude,Longitude,Altitude,TimeZone,DST,TZ, clock_timestamp() from airports;
		`; // add script to use the grpc-connector and load a table
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
		console.log('Removing connection before saving');
		return app.deleteConnection(connectionId)
	})
	.then(() => {
		console.log('Removing script before saving');
		return app.setScript("");
	})
	.then(() => {
		console.log('Saving');
		return app.doSave()
	})
	.then(() => {
		console.log('Fetching Table sample');
		return app.getTableData(-1, 10, true, "Airports")
	})
	.then((tableData) => {
		//Convert table grid into a string using some functional magic
		var tableDataAsString = tableData.map((row) => row.qValue.map((value) => value.qText).reduce((left, right) => left + "\t" + right)).reduce((row1, row2) => row1 + "\n" + row2)
		console.log(tableDataAsString);
	})
	.then(() => session.close())
	.then(() => console.log('Session closed'))
	.catch(err => console.log('Something went wrong :(', err))
