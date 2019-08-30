# Deploy: tell Tilt what YAML to deploy
k8s_yaml('deploy/kubernetes-manifests.yaml')

custom_build(
  'deislabs/smi-adapter-istio',
  'operator-sdk build $EXPECTED_REF',
  ['./pkg', './cmd'],
  tag='latest',
  live_update=[
    sync('pkg', '/go/src/github.com/deislabs/smi-adapter-istio/pkg'),
    run('go install github.com/deislabs/smi-adapter-istio/pkg'),
    restart_container(),
  ]
)
