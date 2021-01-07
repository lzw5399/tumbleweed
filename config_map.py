# -*- coding: utf-8 -*-
import os
import sys
from yaml import load, Loader
import string
import json
from mako.template import Template

PREFIX = 'ENV'

config_map_literal = """apiVersion: v1
kind: ConfigMap
metadata:
  namespace: {}
  name: {}
data:
"""

rbac_literal = """apiVersion: v1
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata: 
    name: {}
roleRef: 
    apiGroup: rbac.authorization.k8s.io
    kind: ClusterRole
    name: admin
subjects: 
- kind: ServiceAccount
  name: default
  namespace: {}
"""

ingress_literal = string.Template("""{{- if .Values.service.enabled -}}
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: $name
  labels:
    app: {{ template "appname" . }}
    chart: "{{ .Chart.Name }}-{{ .Chart.Version| replace "+" "_" }}"
    release: {{ .Release.Name }}
    heritage: {{ .Release.Service }}
  annotations:
    kubernetes.io/tls-acme: "true"
    kubernetes.io/ingress.class: "traefik"
spec:
  rules:
  - host: $host
    http:
      paths:
      - path: /
        backend:
          serviceName: {{ template "fullname" . }}
          servicePort: {{ .Values.service.externalPort }}
{{- end -}}
""")

default_branch_dic = { "develop" : "DEV", "test" : "TEST", "master": "RELEASE", "staging": "STAGING" }

def extend_branch_dic():
    if "BRANCH_DIC" in os.environ:
        brach_dic = json.loads(os.environ['BRANCH_DIC'])
        default_branch_dic.update(brach_dic)
        print(default_branch_dic)

def branch_filter(second_prefix):
    def inner_filter(env):
        words = env.split('_')
        if len(words) > 2 and words[0] == PREFIX and words[1] == second_prefix:
            return True
        return False
    return inner_filter

def trim_prefix(env):
    words = env.split('_')
    return '_'.join(words[2:])

def gen_env_config(p):
    os_env_keys = os.environ.keys()
    env_branch_keys = filter(branch_filter(p), os_env_keys)
    app_env_keys = list(env_branch_keys)
    new_config_env = {trim_prefix(env_key): os.environ[env_key] for env_key in app_env_keys}
    return new_config_env

def gen_configmap():
    current_branch = os.environ['CI_COMMIT_REF_NAME']
    branch_content = current_branch.split('/').pop()
    if default_branch_dic.__contains__(branch_content):
        second_prefixs = default_branch_dic[branch_content]
    config_env = {}

    # 设置通用环境变量
    common_prefix = 'COMMON'
    common_config_env = gen_env_config(common_prefix)
    config_env = {**config_env, **common_config_env}
    
    if default_branch_dic.__contains__(branch_content):
        # 设置分支环境变量
        for p in second_prefixs.split(','):
            new_config_env = gen_env_config(p.strip())
            config_env = {**config_env, **new_config_env}

    fv = open("chart/values.yaml", 'r', encoding='utf8')
    values_yaml = fv.read()
    values_dic = load(values_yaml, Loader=Loader)
    if "env" in values_dic:
        config_env = {**values_dic["env"], **config_env}
    fv.close()

    global config_map_literal
    config_map_literal = config_map_literal.format(os.environ['KUBE_NAMESPACE'], os.environ['CONFIG_MAP_NAME'])

    for config_key, config_value in config_env.items():
        config_map_literal += '  {}: {}\n'.format(config_key, config_value)
    config_map_template = Template(config_map_literal)
    config_map_result = config_map_template.render(**os.environ)
    print(config_map_result)
    fo = open("configmap.yaml", "w")
    fo.write(config_map_result)
    fo.close()

def gen_rbac():
    global rbac_literal
    rbac_literal = rbac_literal.format(os.environ['KUBE_NAMESPACE'], os.environ['KUBE_NAMESPACE'])
    forbac = open("rbac.yaml", "w")
    forbac.write(rbac_literal)
    forbac.close()

def gen_ingress():
    if "URL_PRIFIX" in os.environ:
        print(os.environ['URL_PRIFIX'])
        prifix = os.environ['URL_PRIFIX'].split(',')
        current_branch = os.environ['CI_COMMIT_REF_NAME'].lower()
        for i in range(len(prifix)):
            pure_prifix = prifix[i].replace("'", '').replace('"', '')
            name = '{}-{}'.format(sys.argv[1], i)
            host = '{}-{}'.format(pure_prifix, sys.argv[2])
            if current_branch == 'master':
              host = pure_prifix + '.' + os.environ['AUTO_DEVOPS_DOMAIN']
            global ingress_literal
            blank_str = ingress_literal
            blank_str = blank_str.substitute(vars())
            print(blank_str)
            ingress_path = 'chart/templates/ingress{}.yaml'.format(i)
            fo_ingress = open(ingress_path, "w")
            fo_ingress.write(blank_str)
            fo_ingress.close()

if __name__ == "__main__":
    extend_branch_dic()
    gen_configmap()
    gen_rbac()
    gen_ingress()
