package handlers

import (
	"io/ioutil"
	"net/http"
)

func IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		// w.Write([]byte(indexPage))
		data, _ := ioutil.ReadFile("/home/alshdavid/Development/alshdavid/socks-ssh-manager/daemon/src/httpd/handlers/static/index.html")
		w.Write(data)
	}
}

var indexPage = /* html */ `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Sox5 SSH Manager</title>
  <script src="https://cdn.jsdelivr.net/npm/vue@2"></script>
</head>
<body>
  <div id="app">
    <section>
      <h3>Connected: {{ connected }}</h3>
    </section>
    <section>
      <h3>Client Address</h3>
      <input v-on:input="updateClientAddress(event)" v-model="clientAddress" type="text" placeholder="user@destination.com">
    </section>
    <section>
      <h3>Proxy Strategy</h3>
      <select v-on:input="updateProxyStrategy(event)" v-model="proxyStrategy">
        <option default value="none">Bypass All Traffic</option>
        <option value="all">Proxy All Traffic</option>
        <option value="only-allowed">Proxy Only Listed Domains</option>
        <option value="all-except-denied">Proxy All Except Listed Domains</option>
      </select>
    </section>
    <section>
      <h3>Proxy List</h3>
      <div>
        <label>Domain</label>
        <input v-model="allowedInput" type="text" placeholder="*.google.com"/>
        <button v-on:click="putAllowed()">Add</button>
      </div>
      <div v-for="(allow, i) of allowed">
        <span>{{ allow }}</span>
        <button v-on:click="removeAllowed(i)">Remove</button>
        <br><br>
      </div>
    </section>
    <section>
      <h3>Proxy Bypass List</h3>
      <div>
        <label>Domain</label>
        <input v-model="bypassedInput" type="text" placeholder="*.google.com"/>
        <button v-on:click="putBypassed()">Add</button>
      </div>
      <div v-for="(bypass, i) of bypassed">
        <span>{{ bypass }}</span>
        <button v-on:click="removeBypassed(i)">Remove</button>
        <br><br>
      </div>
    </section>
    <section>
      <h3>Command</h3>
      <code>ssh -D {PORT} {{ clientAddress || 'user@destination.comm' }}</code>
      <button v-on:click="connect()">Connect</button>
    </section>
    <section style="padding-top: 16px;">
      <h3>Logs</h3>
    </section>
    <pre></pre>
  </div>
  <script>
    void async function() {
      const Urls = {
        ProxyList: '/proxy-list',
        ProxyBypassList: '/proxy-bypass-list',
        ProxyStrategy: '/proxy-strategy',
        ClientAddress: '/client-address',
        Connection: '/connection',
      }

      const request = {
        async create(method, url, body) {
          const response = await fetch(url, {
            method,
            headers: body ? { 'Content-Type': 'application/json' } : undefined, 
            body: body ? JSON.stringify(body) : undefined
          })
          if (response.headers.get('Content-Type') === 'application/json') {
            return await response.json()
          }
        },
        put(url, body) { return this.create('PUT', url, body) },
        post(url, body) { return this.create('POST', url, body) },
        delete(url, body) { return this.create('DELETE', url, body) },
        get(url) { return this.create('GET', url) },
      }

      new Vue({
        el: '#app',
        data: {
          connected: (await request.get(Urls.Connection)).status,
          clientAddress: (await request.get(Urls.ClientAddress)).clientAddress,
          proxyStrategy: (await request.get(Urls.ProxyStrategy)).proxyStrategy,
          allowedInput: '',
          allowed: await request.get(Urls.ProxyList),
          bypassedInput: '',
          bypassed: await request.get(Urls.ProxyBypassList),
        },
        methods: {
          updateClientAddress(event) {
            request.put(Urls.ClientAddress, { clientAddress: event.target.value })
          },
          updateProxyStrategy(event) {
            request.put(Urls.ProxyStrategy, { proxyStrategy: event.target.value })
          },
          putAllowed() {
            this.allowed.push(this.allowedInput)
            request.put(Urls.ProxyList, { domain: this.allowedInput })
            this.allowed.sort()
            this.allowedInput = ''
          },
          removeAllowed(i) {
            const domain = this.allowed.splice(i, 1)[0]
            request.delete(Urls.ProxyList, { domain })
          },
          putBypassed() {
            this.bypassed.push(this.bypassedInput)
            request.put(Urls.ProxyBypassList, { domain: this.bypassedInput })
            this.bypassedInput = ''
          },
          removeBypassed(i) {
            const domain = this.bypassed.splice(i, 1)[0]
            request.delete(Urls.ProxyBypassList, { domain })
          },
          async connect() {
            await request.post(Urls.Connection, { action: 'CONNECT' })
            while (this.connected === false) {
              this.connected = (await request.get(Urls.Connection)).status
              await new Promise(res => setTimeout(res, 250))
            }
          },
          async disconnect() {
            await request.post(Urls.Connection, { action: 'DISCONNECT' })
            while (this.connected === true) {
              this.connected = (await request.get(Urls.Connection)).status
              await new Promise(res => setTimeout(res, 250))
            }
          },
        }
      })
    }()
  </script>
  <style>
    * {
      margin: 0;
    }

    body, #app {
      height: 100vh;
    }

    #app {
      padding-top: 16px;
    }

    h3, input, select, label, code {
      margin-bottom: 16px;
      display: block;
    }

    section {
      margin-left: 16px;
    }

    input {
      display: inline-block;
    }

    label {
      margin-bottom: 4px;
    }

    pre {
      border-top: 1px solid black;
    }
  </style>
</body>
</html>
`
