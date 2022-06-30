allow_k8s_contexts('wego-dev')

image_repository = os.getenv('IMAGE_REPO', 'localhost:5001/weaveworks/progressive-delivery')

load('ext://restart_process', 'docker_build_with_restart')

docker_build(
  image_repository,
  '.',
  dockerfile="tools/tilt/dev-server.dockerfile",
)

k8s_yaml('tools/tilt/service-account.yaml')
k8s_yaml('tools/tilt/role.yaml')
k8s_yaml('tools/tilt/app.yaml')

k8s_resource('progressive-delivery-server', port_forwards=9002)
