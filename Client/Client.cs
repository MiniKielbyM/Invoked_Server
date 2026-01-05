using System;
using System.Net;
using System.Net.Sockets;
using System.Text;

class Program
{
    static void Main()
    {
        try
        {
            // Connect to server
            using var client = new UdpClient();
            var serverEndpoint = new IPEndPoint(IPAddress.Parse("127.0.0.1"), 8080);
            client.Connect(serverEndpoint);

            // Send a message
            SendServerMessage(client, new object[] { "header1", "header2" }, "Hello, UDP!");

            // Wait for reply with timeout
            client.Client.ReceiveTimeout = 5000; // 5 seconds
            var remoteEndpoint = new IPEndPoint(IPAddress.Any, 0);
            byte[] response = client.Receive(ref remoteEndpoint);

            Console.WriteLine($"Server says: {Encoding.UTF8.GetString(response)}");
        }
        catch (SocketException ex)
        {
            Console.WriteLine($"Socket error: {ex.Message}");
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Error: {ex.Message}");
        }
    }
    static void HandleError(string message)
    {
        Console.WriteLine($"Error: {message}");
        Environment.Exit(1);
    }

    static void SendServerMessage(UdpClient client, Object[] head, string message)
    {
        byte[] data = Encoding.UTF8.GetBytes(message);
        byte[] header = Encoding.UTF8.GetBytes("||HEADER.START||" + string.Join(",", head) + "||HEADER.END||");
        byte[] combined = new byte[header.Length + 1 +data.Length]; // +1 for delimiter
        Buffer.BlockCopy(header, 0, combined, 0, header.Length);
        Buffer.BlockCopy(data, 0, combined, header.Length + 1, data.Length);
        client.Send(combined, combined.Length);
    }
}
