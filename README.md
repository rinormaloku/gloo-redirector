# Gloo Redirector

Generate redirection configuration from a CSV file for Gloo Edge and Gloo Mesh.

**Install Gloo Redirector** by executing:
```bash
go install github.com/rinormaloku/gloo-redirector@latest
```

And proceed with the examples below.

## How to generate redirections

Write all the redirections in a CSV file. For example, to redirect traffic from the **initial_location** `https://solo.io/docs/a` to the **redirect_location** `https://docs.solo.io/a` with a **redirect_code** of `301` you'd use the following sequence: `initial_location,redirect_location,redirect_code`. 

Thus, for the above example we'd write the following file (keep adding further entries in new lines):
```
cat <<EOF > /tmp/redirections.csv
https://solo.io/docs/a,https://docs.solo.io/a,301
EOF
```

Then use the CSV as the source in the command below:
```
gloo-redirector edge generate --source /tmp/redirections.csv
```

This produces the virtual services to configure the gateway proxies for redirecting traffic, as shown below:

```yaml
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: redirect-solo-io
  namespace: gloo-system
spec:
  virtualHost:
    domains:
      - solo.io
    routes:
      - matchers:
          - exact: /docs/a
        redirectAction:
          hostRedirect: docs.solo.io
          pathRedirect: /a
          responseCode: MOVED_PERMANENTLY
```

### Customizing the default template

The above output is based on a default template. The default one however, might not fit with your use-case.
You can print the template with the command below:
```bash
$ gloo-redirector edge print-template 

apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: redirect-{{ .ResourceName }}
  namespace: gloo-system
spec:
  virtualHost:
    domains:
      - {{  .Host }}
    routes:
    {{- range $matcher := .Matchers }}
      - matchers:
          - exact: {{ $matcher.ExactPath }}
        redirectAction:
          hostRedirect: {{ $matcher.HostRewrite }}
          pathRedirect: {{ $matcher.PathRewrite }}
          responseCode: {{ $matcher.RedirectCode }}
    {{- end }}
```

Redirect the output to a file (e.g. `/tmp/template.yaml`) and modify it to your liking.
Then you can execute the generate command with your template:

```bash
gloo-redirector edge generate --source /tmp/redirections.csv --template /tmp/template.yaml
```

**TIP**: You can pipe the output to kubectl:
```bash
gloo-redirector edge generate --source /tmp/redirections.csv | kubectl apply -f - 
```

## Gloo Redirector Help
```bash
$ gloo-redirector --help

Gloo Redirector generates 3xx redirection configuration for either Gloo Edge and Gloo Mesh.
Examples:
  # Generate Gloo Mesh redirection configuration using a file as a source with the default template
  gloo-redirector mesh generate --source /tmp/redirections.csv

  # Generate Gloo Edge redirection configuration using a file as a source with the default template
  gloo-redirector edge generate --source /tmp/redirections.csv

Usage:
  gloo-redirector [command]

Available Commands:
  edge        Gloo Edge commands
  mesh        Gloo Mesh commands

Flags:
  -h, --help   help for gloo-redirector

Use "gloo-redirector [command] --help" for more information about a command.
```

**NOTE:** [istio-redirector](https://github.com/blablacar/istio-redirector) was the starting point for this CLI.
