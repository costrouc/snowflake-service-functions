USE ROLE NATIVE_APP_PUBLISHER;

CREATE SERVICE IF NOT EXISTS CNO_SERVICE_FUNCTIONS.DATA.service_function
  IN COMPUTE POOL CNO_SERVICE_FUNCTIONS
  FROM SPECIFICATION $$
spec:
  containers:
  - name: servicefunctions
    image: /cno_service_functions/data/image_repository/service-functions:latest
  endpoints:
  - name: http
    port: 8080
  $$
  AUTO_RESUME = TRUE
  MIN_INSTANCES = 1
  MAX_INSTANCES = 1;

SHOW SERVICE CONTAINERS IN SERVICE CNO_SERVICE_FUNCTIONS.DATA.service_function;

CREATE FUNCTION IF NOT EXISTS CNO_SERVICE_FUNCTIONS.DATA.DEBUG (text varchar)
   RETURNS varchar
   SERVICE=CNO_SERVICE_FUNCTIONS.DATA.service_function
   ENDPOINT=http
   AS '/rpc';

SELECT SYSTEM$GET_SERVICE_STATUS('CNO_SERVICE_FUNCTIONS.DATA.service_function', 600);
