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

            // Set timeout for receiving responses
            client.Client.ReceiveTimeout = 5000;

            Console.WriteLine("Connected to server. Type 'exit' to quit.");

            while (true)
            {
                // Get user input
                Console.Write("\nEnter message: ");
                string? userInput = Console.ReadLine();

                if (userInput?.ToLower() == "exit")
                {
                    break;
                }

                if (string.IsNullOrEmpty(userInput))
                {
                    continue;
                }

                if (userInput.ToLower() == "start game")
                {
                    StartGame(writer);
                    continue;
                }

                // Send message to server
                SendServerMessage(writer, new object[] { "header1", "header2" }, userInput);

                try
                {
                    // Read response from server
                    string? response = reader.ReadLine();
                    if (response != null)
                    {
                        Console.WriteLine($"Server says: {response}");
                    }
                }
                catch (SocketException)
                {
                    Console.WriteLine("No response from server (timeout)");
                }
            }

            Console.WriteLine("Disconnected from server.");
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
        string header = string.Join("||HEADER.SEP||", head) + "||HEADER.END||";
        writer.WriteLine(header + message);
    }

    static void StartGame(StreamWriter writer)
    {
        SendServerMessage(writer, new object[] { "game", "start" }, "");
    }
}
