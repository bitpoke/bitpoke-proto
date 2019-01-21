import fileinput
import yaml


metadata = open('hack/chart-metadata.yaml', 'r')
parsed_metadata = yaml.safe_load(metadata)
metadata.close()

document = open('config/webhook/webhook.yaml', 'r')
loaded_webhooks = list(yaml.safe_load_all(document))
document.close()

for num, doc in enumerate(loaded_webhooks):
    doc['metadata']['labels'] = parsed_metadata['metadata']['labels']
    doc['metadata'].pop('creationTimestamp', None)

    if num == 0:
        # secret
        doc['data']['ca-cert.pem'] = '{{ $cert.Cert | b64enc }}'
        doc['data']['ca-key.pem'] = '{{ $cert.Key | b64enc }}'
        doc['data']['cert.pem'] = '{{ $cert.Cert | b64enc }}'
        doc['data']['key.pem'] = '{{ $cert.Key | b64enc }}'
        doc['metadata']['name'] = '{{ include "dashboard.webhook.secretName" . }}'
        doc['metadata']['namespace'] = '{{ .Release.Namespace }}'
        doc['metadata']['annotations'] = doc['metadata'].get('annotations', {})
        doc['metadata']['annotations']['helm.sh/hook'] = 'pre-install'
        doc['metadata']['annotations']['helm.sh/hook-delete-policy'] = 'before-hook-creation'
    elif num == 1:
        # service
        doc['metadata']['name'] = '{{ include "dashboard.fullname" . }}-webhook'
        doc['metadata']['namespace'] = '{{ .Release.Namespace }}'
        doc['spec']['type'] = '{{ .Values.webhook.service.type }}'
        doc['spec']['ports'][0]['port'] = '{{ .Values.webhook.service.port }}'
        doc['spec']['ports'][0]['targetPort'] = '{{ .Values.webhook.service.targetPort }}'
        doc['spec']['ports'][0]['protocol'] = 'TCP'
        doc['spec']['ports'][0]['name'] = 'http'
        doc['spec']['selector']['app'] = '{{ include "dashboard.name" . }}'
        doc['spec']['selector']['release'] = '{{ .Release.Name }}'
    elif num == 2:
        # webhooks
        for webhook in doc['webhooks']:
            webhook['clientConfig']['caBundle'] = '{{ $cert.Cert | b64enc }}'
            webhook['clientConfig']['service']['name'] = '{{ include "dashboard.fullname" . }}-webhook'
            webhook['clientConfig']['service']['namespace'] = '{{ .Release.Namespace }}'

parsed_webhooks = open("chart/dashboard/templates/webhook.yaml", "w")
yaml.dump_all(loaded_webhooks, parsed_webhooks, default_flow_style=False)
parsed_webhooks.close()

# overrides
with fileinput.FileInput("chart/dashboard/templates/webhook.yaml", inplace=True) as file:
    for line in file:
        line = line.replace("'{{ .Values.webhook.service.port }}'", "{{ .Values.webhook.service.port }}")
        line = line.replace("'{{ .Values.webhook.service.targetPort }}'", "{{ .Values.webhook.service.targetPort }}")
        print(line, end='')
