import { createClient } from "@connectrpc/connect";
import { createGrpcWebTransport} from "@connectrpc/connect-web";
import { ReaderService } from "../api/reader_pb";


const transport = createGrpcWebTransport({
  baseUrl: "http://localhost:50051",
  useBinaryFormat:false
});

// Here we make the client itself, combining the service
// definition with the transport.
export const client = createClient(ReaderService, transport);