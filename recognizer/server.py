from concurrent import futures
import logging


import grpc
from grpc_reflection.v1alpha import reflection
import recognizer_pb2_grpc as recognizer_pb2_grpc
import recognizer_pb2 as recognizer_pb2


class RecognizerServicer(recognizer_pb2_grpc.RecognizerServicer):
    def __init__(self):
        pass

    def RecognizeFormula(self, request, context):
        return recognizer_pb2.RecognizedResult(result=f'formula{request}')

    def RecognizeMixed(self, request, context):
        logging.log(logging.INFO, request)
        return recognizer_pb2.RecognizedResult(result=f'mixed{request}')

    def RecognizeText(self, request, context):
        logging.log(logging.INFO, request)
        return recognizer_pb2.RecognizedResult(result=f'text{request}')


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    recognizer_pb2_grpc.add_RecognizerServicer_to_server(
        RecognizerServicer(), server
    )
    SERVICE_NAMES = (
        recognizer_pb2.DESCRIPTOR.services_by_name['Recognizer'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)
    server.add_insecure_port("[::]:50051")
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig()
    serve()
