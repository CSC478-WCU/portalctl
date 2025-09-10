# USAGE

```
portalctl <command> [global flags] [command flags]
```

**Global flags (apply to all commands):**

```
  -server string   XML-RPC server (default "boss.emulab.net")
  -port int        XML-RPC port (default 3069)
  -path string     Server path (default "/usr/testbed")
  -cert string     Client cert PEM (or combined PEM)
  -key string      Client key PEM (defaults to -cert)
  -cacert string   CA cert PEM (required if -verify)
  -verify          Verify server cert using CA (needs to be removed, certs are self-signed and throw errors)
  -timeout int     HTTP timeout seconds (default 900)
```

**Commands:**

```
  start       -project <pid[,gid]> -name <eid>
              [--bindings JSON | --bindings-file PATH | --param k=v ...]
              [--refspec url:ref] [--aggregate URN] [--site SITESTR]
              [--duration H] [--start T] [--stop T] [--sshpubkey STR]
              <profile-uuid-or-pid,name>

  status      [-j] [-k] [-r] <experiment-uuid-or-pid,name>

  modify      [--bindings JSON | --bindings-file PATH | --param k=v ...]
              <experiment-uuid-or-pid,name>

  terminate   <experiment-uuid-or-pid,name>

  extend      [-m reason] <experiment-uuid-or-pid,name> <hours>

  manifests   <experiment-uuid-or-pid,name>

  reboot      [-f] <experiment-uuid-or-pid,name> node [node ...]

  connect     <src-exp> <src-lan> <dst-exp> <dst-lan>

  disconnect  <experiment-uuid-or-pid,name> <src-lan>
```

**Notes:**

- `--bindings` must be a **JSON object** mapping **parameter names → string values**.If a profile expects JSON-typed parameters, pass them as **stringified JSON**.
- `--param k=v` is a convenience to build the `bindings` object on the fly (repeatable).
- For `status`, `-j` makes the server return a JSON payload string.

**Examples:**

_Start the terraform-profile experiment (with xvenm on rawpc using [/examples/xenvm-on-rawpc.json](/examples/xenvm-on-rawpc.json)):_

```bash
 ./portalctl start \
  -server boss.emulab.net \
  -path /usr/testbed \
  -cert ./cert.pem \
  -project cloud-edu \
  -name xenvm-on-rawpc \
  -spec-file ./examples/xvenm-on-rawpc.json cloud-edu,terraform-profile
```

_Start the default [emulab-ops,k8s]([url](https://www.cloudlab.us/show-profile.php?uuid=79d36573-a099-11ea-b1eb-e4434b2381fc)) profile with default params_ 
```bash
./portalctl start -server boss.emulab.net -path /usr/testbed   -cert ./cert.pem   -project cloud-edu -name tf-demo emulab-ops,k8s
```

_Check status (JSON payload from server):_

```bash
./portalctl status -j cloud-edu,tf-demo
```

_Modify parameters on a running experiment (from a file):_

```bash
./portalctl modify --bindings-file ./bindings.json cloud-edu,tf-demo
```

_Terminate an experiment:_

```bash
./portalctl terminate cloud-edu,tf-demo
```

# Build tips

- **Requirements:** Go 1.22+
- **Clone & build:**

  ```bash
  git clone https://github.com/CSC478-WCU/portalctl.git
  cd portalctl
  go build -o portalctl ./cmd/portalctl
  ```
- **Windows build:**

  ```powershell
  go build -o portalctl.exe .\cmd\portalctl
  ```
- **Run from anywhere:** add the compiled binary’s folder to your `PATH`.
- **TLS files:**

  - Use `-cert` and `-key` to point to your client certificate and key (PEM).If you have a single combined PEM, pass it to both flags (or only `-cert` when the key is included).
  - To verify the server certificate, pass `-verify -cacert <CA.pem>`.
- **Troubleshooting TLS:**

  - If you see errors like *“invalid authority info access”*, ensure you’re using separate `-cert` and `-key` PEMs (or a combined PEM) 
  - If your cert/key are encrypted, decrypt to plain PEMs before use.
