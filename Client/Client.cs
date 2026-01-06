using System;
using System.IO;
using System.Net.Sockets;
using System.Text;

class Client
{
    static void Main()
    {
        try
        {
            // Connect to server
            using var client = new TcpClient("127.0.0.1", 8080);
            using var stream = client.GetStream();
            using var reader = new StreamReader(stream, Encoding.UTF8);
            using var writer = new StreamWriter(stream, Encoding.UTF8) { AutoFlush = true };

            // Send a message
            SendServerMessage(writer, new object[] { "header1", "header2" }, "Hello, TCP!");
            client.ReceiveTimeout = 5000; // 5 seconds
            string? response = reader.ReadLine();
            Console.WriteLine($"Server says: {response}");
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

    static void SendServerMessage(StreamWriter writer, Object[] head, string message)
    {
        string header = "||HEADER.START||" + string.Join(",", head) + "||HEADER.END||";
        writer.WriteLine(header + message);
    }
}
