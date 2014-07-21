<h1>IB</h1>

<p>A golang package for the <a href="http://www.interactivebrokers.com">Interactive Brokers</a> API.</p>

<h2>Package Details</h2>

<p>
	<a href="http://www.interactivebrokers.com">IB Gateway</a> implements a streaming message protocol consisting of null byte (\000) delimited strings. 
	As a result, the data that is sent and received follows a well-defined pattern that must be observed.
	Each request and response consists of an integral code, a version number, and the relevant data. See the IB Gateway source for details on the protocol.
</p>
<h3>The <span style="font-family: monospace;">Broker</span> Type</h3>
<p>
	The main goal of this package is to encapsulate the functionality required for reading and writing messages that are required by the IB Gateway application.
	The <span style="font-family: monospace">broker.go</span> file contains a base type called <span style="font-family: monospace">Broker</span>, which can be inherited by other types.
	The <span style="font-family: monospace">Broker</span> type implements functions that are used across the package by various "brokers."
	For example, the <span style="font-family: monospace">marketdata.go</span> file houses a <span style="font-family: monospace">MarketDataBroker</span> type that inherits the <span style="font-family: monospace">Broker</span> type.
	The <span style="font-family: monospace">MarketDataBroker</span> definition includes an anonymous <span style="font-family: monospace">Broker</span> field that accepts an empty <span style="font-family: monospace">Broker</span> struct when the <span style="font-family: monospace">NewMarketDataBroker</span> function is called.
</p>
<h3>Brokers</h3>
<p>
	As mentioned above, the IB package consists of several "brokers" that inherit base functionality from the <span "font-family: monospace;">Broker</span> type.
	These brokers correspond to the functionalitity implemented in the IB Gateway C++ API.
	<table>
		<thead>
			<tr style="font-weight: 600;">
				<td><a href="http://www.interactivebrokers.com">IB Gateway Documentation</a></td>
				<td>IB Package Type</td>
				<td>IB Package File</td>
			</tr>
		</thead>
		<tbody>
			<tr>
				<td>Historical Data</td>
				<td>HistoricalDataBroker</td>
				<td>historicaldata.go</td>
			</tr>
			<tr>
				<td>Market Data</td>
				<td>MarketDataBroker</td>
				<td>marketdata.go</td>
			</tr>
			<tr>
				<td>Contract Details</td>
				<td>ContractDetailsBroker</td>
				<td>contractdetails.go</td>
			</tr>
			<tr>
				<td>Account Value</td>
				<td>AccountDataBroker</td>
				<td>accountdata.go</td>
			</tr>
			<tr>
				<td>Orders</td>
				<td>OrderBroker</td>
				<td>order.go</td>
			</tr>
		</tbody>
	</table>
	
	<h4>The "New" Function</h4>
	Each broker has a "new" function that returns an initialized broker.
	<pre>
	mktData := NewMarketDataBroker()
	</pre>

	<h4>The <span style="font-family: monospace;">Connect</span> Function</h4>
	The base <span style="font-family: monospace;">Broker</span> type has a <span style="font-family: monospace;">Connect</span> function that connects the broker to the TCP/IP socket that the IB Gateway application is serving. The <span style="font-family: monospace;">Broker</span> type also has a <span style="font-family: monospace;">Disconnect</span> function.
	<pre>
	err = mktData.Connect()
	defer mktData.Disconnect()
	</pre>

	<h4>"Listening" for Response Data</h4>
	The <span style="font-family: monospace;">Connect</span> function must be called before attempting to communicate with the IB Gateway application.
	Once the broker is connected, a client application can invoke the <span style="font-family: monospace;">Listen</span> function in a goroutine.
	The <span style="font-family: monospace;">Listen</span> function accepts an anonymous callback function whose signature is parameterless and returns no values.
	The callback function has access to the broker and is therefore able to respond to messages sent on channels that are exposed in the broker.
	<pre>
	go mktData.Listen(func () {
			for {
				select {
					case p := <lt;-m.TickPriceChan:
						ib.Log.Print("price", p)
					case s := <lt;-m.TickSizeChan:
						ib.Log.Print("size", s)
					case o := <lt;-m.TickOptCompChan:
						ib.Log.Print("optcomp", o)
					case g := <lt;-m.TickGenericChan:
						ib.Log.Print("generic", g)
					case t := <lt;-m.TickStringChan:
						ib.Log.Print("string", t)
					case e := <lt;-m.TickEFPChan:
						ib.Log.Print("efp", e)
					case d := <lt;-m.MarketDataTypeChan:
						ib.Log.Print("datatype", d)
				}
			}
		})
	</pre>
	<p><span style="color: orange;">There are probably better, more flexible ways of implementing the listener functionality. One possibility is allowing/forcing a client application to kick off a goroutine before invoking the <span style="font-family: monospace;">Listen</span> function. The end result is still 2 goroutines for each listener. Although, if you forget to start a goroutine that responds to data sent on the data channels, you won't hear about it on compilation.</span></p>

	<h4>Data Channels</h4>
	The brokers have channels over which reponse data objects are sent. 
	Each broker exposes a different set of channels for communicating responses from the IB Gateway application.
	As shown above, the callback that is supplied to the <span style="font-family: monospace;">Listen</span> function is able to access the data channels in a goroutine.
</p>
<h3>Configuration, Globals, Requests and Responses</h3>
<p>
	<h4><span style="font-family: monospace;">conf.go</span></h4>
	<p>Configuration variables i.e. host, port, version, etc.</p>
	<h4><span style="font-family: monospace;">globals.go</span></h4>
	<p>Useful global variables.</p>
	<h4><span style="font-family: monospace;">requestcodes.go</span></h4>
	<p>Request codes and version variables. Initialized in a call to <span style="font-family: monospace;">init</span>.</p>
	<h4><span style="font-family: monospace;">responsecodes.go</span></h4>
	<p>Response codes. Initialized in a call to <span style="font-family: monospace;">init</span>.</p>
	<p><span style="color: orange;">Not using .ini file here. Just a personal aversion to library/package/framework dependencies. Not married to this style though.</span></p>
</p>
<h2>Running IB Gateway</h2>
<p>
	Running IB Gateway on a server poses a couple of small obstacles.
</p>
<p>
	First of all, the IB Gateway GUI is inseparable from the functionality it provides. Therefore, the GUI must be running to communicate with the brokerage. Since the GUI is somewhat minimal, X11 forwarding is not a major issue. But we must divert the GUI and keep it running while the algorithms run. To accomplish this, install <a href="http://www.xpra.org">xpra</a>. Described as "screen for GUI applications," xpra will let us start IB Gateway, login, and subsequently detach while keeping the application running in an xpra session. 
	<p>First, the user must log in to the server with X11 forwarding over ssh.<p>
	<p>To start an xpra server, use <span style="font-family: monospace;">xpra start :0</span>.</p>
	<p>Once the xpra server is running, start a screen session on display 0: <span style="font-family: monospace;">DISPLAY=:0 screen</span>.</p>
	<p>In the screen session, start IB Gateway in the background: <span style="font-family: monospace;">./startup.sh &</span>
	<p>Detach from screen via &lt;CTL-d&gt;.</p>
	<p>Attach to the xpra server with display 0: <span style="font-family: monospace;">xpra attach :0</span></p>
	<p>Attaching to the xpra server will start X11 forwarding of the IB Gateway GUI and the standard log in steps can be completed.</p>
	<p>Once logged in, the user can safely detach from xpra with &lt;CTL-c&gt;. IB Gateway will continue to run but the GUI will be diverted.</p>
	<p>Run algos...</p>
</p> 
