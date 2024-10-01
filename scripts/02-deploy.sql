USE ROLE ACCOUNTADMIN;

CREATE SERVICE IF NOT EXISTS CNO_SERVICE_FUNCTIONS.DATA.service_function
  IN COMPUTE POOL CNO_SERVICE_FUNCTIONS
  FROM SPECIFICATION $$
spec:
  containers:
  - name: service-function
    image: /cno_service_functions/data/image_registry/service-functions:latest
  $$
  AUTO_RESUME = TRUE;
  MIN_INSTANCES = 1
  MAX_INSTANCES = 2;
