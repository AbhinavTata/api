import datetime
import grpc
from google.protobuf.json_format import MessageToJson
import json

from . import araali_api_service_pb2
from . import araali_api_service_pb2_grpc
from . import utils

def init_assetsreq(zone, app, ago):
    req = araali_api_service_pb2.ListAssetsRequest()
    
    if utils.cfg["tenant"]: req.tenant.id = utils.cfg["tenant"]
    
    if zone: req.zone = zone
    if app: req.app = app

    req.filter.list_active_vm = True
    req.filter.list_inactive_vm = False
    req.filter.list_active_container = True
    req.filter.list_inactive_container = False

    if not ago:
        ago = "days=1"
    ago = dict([utils.make_map(a) for a in ago.split(",")])
    req.time.start_time.FromDatetime(datetime.datetime.now() - datetime.timedelta(**ago))
    req.time.end_time.FromDatetime(datetime.datetime.now())
    
    print(req)

    return req

def init_linksreq(zone, app, svc, ago):
    req = araali_api_service_pb2.ListLinksRequest()
    
    if utils.cfg["tenant"]: req.tenant.id = utils.cfg["tenant"]
    if zone: req.zone = zone
    if app: req.app = app
    if svc: req.service = svc

    if not ago:
        ago = "days=1"
    ago = dict([utils.make_map(a) for a in ago.split(",")])
    req.time.start_time.FromDatetime(datetime.datetime.now() - datetime.timedelta(**ago))
    req.time.end_time.FromDatetime(datetime.datetime.now())
    
    print(req)

    return req

def init_alertreq(count, ago):
    req = araali_api_service_pb2.ListAlertsRequest()
    
    if utils.cfg["tenant"]: req.tenant.id = utils.cfg["tenant"]
    
    req.filter.open_alerts = True
    req.filter.closed_alerts = False
    req.filter.list_all_alerts = False
    req.filter.perimeter_ingress = True
    req.filter.perimeter_egress = True
    req.filter.home_non_araali_ingress = True
    req.filter.home_non_araali_egress = True
    req.filter.araali_to_araali = True
    if not ago:
        ago = "infinite"
    if ago != "infinite":
        ago = dict([utils.make_map(a) for a in ago.split(",")])
        req.filter.time.start_time.FromDatetime(datetime.datetime.now() - datetime.timedelta(**ago))
        req.filter.time.end_time.FromDatetime(datetime.datetime.now())
    
    if count is None: count = 1000
    req.count = count
    print(req)

    return req

def init_insightsreq(zone):
    req = araali_api_service_pb2.ListInsightsRequest()
    
    if utils.cfg["tenant"]: req.tenant.id = utils.cfg["tenant"]
    if zone: req.zone = zone
    
    print(req)

    return req

class API:
    def __init__(self):
        if utils.cfg["token"] is None:
            raise Exception("*** ARAALI_API_TOKEN environment variable is not set")

        class GrpcAuth(grpc.AuthMetadataPlugin):
            def __init__(self, key):
                self._key = key

            def __call__(self, context, callback):
                callback((('authorization', self._key),), None)

        token_cred = grpc.metadata_call_credentials(GrpcAuth("Bearer %s" % utils.cfg["token"]))
        channel_cred = grpc.ssl_channel_credentials()
        creds = grpc.composite_channel_credentials(channel_cred, token_cred)

        backend = utils.cfg["backend"]
        if not backend:
            backend = "prod"
        channel = grpc.secure_channel('api-%s.aws.araalinetworks.com:443' % (backend), creds)

        self.stub = araali_api_service_pb2_grpc.AraaliAPIStub(channel)

    def get_alerts(self, count=None, ago=None):
        """Fetches alerts
            Usage: alerts, next_page, status = api.get_alerts()
        """
        resp = self.stub.listAlerts(init_alertreq(count, ago))
        if resp.response.code != 0:
            print("*** Error fetching alerts:", resp.response.message)
        return ([json.loads(MessageToJson(a)) for a in resp.links], resp.paging_token, resp.response.code)

    def get_assets(self, zone=None, app=None, ago=None):
        """Fetches assets
            Usage: assets, status = api.get_assets()
        """
        resp = self.stub.listAssets(init_assetsreq(zone, app, ago))
        if resp.response.code != 0:
            print("*** Error fetching alerts:", resp.response.message)
        return ([json.loads(MessageToJson(a)) for a in resp.assets], resp.response.code)

    def get_links(self, zone=None, app=None, svc=None, ago=None):
        """Fetches links
            Usage: links, status = api.get_links()
        """
        resp = self.stub.listLinks(init_linksreq(zone, app, svc, ago))
        if resp.response.code != 0:
            print("*** Error fetching alerts:", resp.response.message)
        return ([json.loads(MessageToJson(a)) for a in resp.links], resp.response.code)

    def get_insights(self, zone=None):
        """Fetches insights
            Usage: insights, status = api.get_insights()
        """
        resp = self.stub.listInsights(init_insightsreq(zone))
        if resp.response.code != 0:
            print("*** Error fetching alerts:", resp.response.message)
        return ([json.loads(MessageToJson(a)) for a in resp.insights], resp.response.code)
