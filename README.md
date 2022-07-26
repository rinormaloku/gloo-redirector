# Gloo Redirector

Quickly generate RouteTables and VirtualGateways for redirecting traffic

> NOTE: This is a fork of istio-redirector, and slimmed down to the CLI. 

Install by executing:
```bash
go install github.com/rinormaloku/gloo-redirector@latest
```

Now to execute commands use the binary `gloo-redirector`.

## How to generate redirections?

Write all the redirections in a file formatted as comma seperated values. E.g.:
```
cat <<EOF > /tmp/redirections.csv
https://solo.io/docs/a,https://docs.solo.io/a,301
https://solo.io/docs/b,https://docs.solo.io,301
https://solo.io/landingpage/a,https://landing.solo.io/a,303
https://test.solo.io/,https://staging.solo.io/,303
EOF
```

Then execute generate and use the csv as a source:
```
gloo-redirector generate --source /tmp/redirections.csv
```

This produces a series of VirtualGateways and RouteTables, to configure the gateway(s) for redirection.

### Customizing the Gateways and the RouteTables

By default, you will create the following virtualgateway and route table definition shown below. 
Understandably, you have to change this to match your Gateways.

Do so by simply writing this to a location and modifying it to your liking. E.g.
```yaml
cat <<EOF > /tmp/template.yaml
apiVersion: networking.gloo.solo.io/v2
kind: VirtualGateway
metadata:
  name: north-south-{{ .RouteTableName }}
  namespace: istio-gateways
  labels:
    gloo: redirector
spec:
  workloads:
    - selector:
        labels:
          istio: ingressgateway
        cluster: cluster1
  listeners:
    - http: {}
      port:
        number: 80
      allowedRouteTables:
        - host: {{ .Host }}
---
apiVersion: networking.gloo.solo.io/v2
kind: RouteTable
metadata:
  name: redirect-{{ .RouteTableName }}
  namespace: istio-gateways
  labels:
    gloo: redirector
spec:
  hosts:
    - {{  .Host }}
  virtualGateways:
    - name: north-south-{{ .RouteTableName }}
      namespace: istio-gateways
      cluster: cluster1
  http:
    {{- range $matcher := .Matchers }}
    - matchers:
        - uri:
            exact: {{ $matcher.ExactPath }}
      redirect:
        hostRedirect: {{ $matcher.HostRewrite }}
        pathRedirect: {{ $matcher.PathRewrite }}
        responseCode: {{ $matcher.RedirectCode }}
    {{- end }}
EOF
```

> **NOTE:** You mainly will want to change the VirtualGateway selector and namespace.

Then you can execute the generate command with your template:

```bash
gloo-redirector generate --source /tmp/redirections.csv --template /tmp/template.yaml
```

**TIP**: You can pipe the output to kubectl:
```bash
gloo-redirector generate --source /tmp/redirections.csv | kubectl apply -f - 
```
