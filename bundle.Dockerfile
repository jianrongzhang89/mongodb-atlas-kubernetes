FROM scratch
<<<<<<< HEAD
<<<<<<< HEAD

LABEL com.redhat.openshift.versions="v4.5-v4.7"
=======
=======

# RH OLM annotations
>>>>>>> dbaas bundle build
LABEL com.redhat.openshift.versions="v4.5-v4.8"
>>>>>>> Integrate Atlas Operator with Red Hat DBaaS
LABEL com.redhat.delivery.backport=true
LABEL com.redhat.delivery.operator.bundle=true

# Core bundle labels.
LABEL operators.operatorframework.io.bundle.mediatype.v1=registry+v1
LABEL operators.operatorframework.io.bundle.manifests.v1=manifests/
LABEL operators.operatorframework.io.bundle.metadata.v1=metadata/
LABEL operators.operatorframework.io.bundle.package.v1=mongodb-atlas-kubernetes
LABEL operators.operatorframework.io.bundle.channels.v1=beta
LABEL operators.operatorframework.io.bundle.channel.default.v1=beta
<<<<<<< HEAD
<<<<<<< HEAD
LABEL operators.operatorframework.io.metrics.builder=operator-sdk-v1.16.0
LABEL operators.operatorframework.io.metrics.mediatype.v1=metrics+v1
LABEL operators.operatorframework.io.metrics.project_layout=go.kubebuilder.io/v2
=======
LABEL operators.operatorframework.io.metrics.builder=operator-sdk-v1.7.1+git
=======
LABEL operators.operatorframework.io.metrics.builder=operator-sdk-v1.10.0+git
>>>>>>> dbaas bundle build
LABEL operators.operatorframework.io.metrics.mediatype.v1=metrics+v1
LABEL operators.operatorframework.io.metrics.project_layout=go.kubebuilder.io/v3
>>>>>>> Integrate Atlas Operator with Red Hat DBaaS

# Labels for testing.
LABEL operators.operatorframework.io.test.mediatype.v1=scorecard+v1
LABEL operators.operatorframework.io.test.config.v1=tests/scorecard/

# Copy files to locations specified by labels.
COPY bundle/manifests /manifests/
COPY bundle/metadata /metadata/
COPY bundle/tests/scorecard /tests/scorecard/
