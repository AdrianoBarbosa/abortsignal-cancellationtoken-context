using System.Threading.Channels;

namespace Cancelamento.Demos;

/// <summary>
/// Example 5: a wait with no way out, the .NET version of a goroutine leak.
///
/// Three readers wait on a Channel that nobody writes to or completes.
/// Without a token, ReadAsync never returns and the Tasks stay pending for as
/// long as the process lives. With a token they all leave when the
/// cancellation arrives.
/// </summary>
public static class EsperaSemSaida
{
    public static async Task RodarAsync()
    {
        await SemTokenAsync();
        Console.WriteLine();
        await ComTokenAsync();
    }

    private static async Task SemTokenAsync()
    {
        Console.WriteLine("-- ReadAsync sem token: as tarefas ficam esperando --");
        var canal = Channel.CreateUnbounded<int>();

        var leitores = Enumerable.Range(1, 3)
            .Select(_ => canal.Reader.ReadAsync().AsTask()) // nobody will ever write
            .ToArray();

        await Task.Delay(300);
        Console.WriteLine($"  tarefas ainda esperando: {leitores.Count(t => !t.IsCompleted)} de {leitores.Length}");
    }

    private static async Task ComTokenAsync()
    {
        Console.WriteLine("-- ReadAsync com token: todas saem --");
        var relogio = Apoio.Cronometrar();
        var canal = Channel.CreateUnbounded<int>();

        using var cts = new CancellationTokenSource(TimeSpan.FromMilliseconds(300));

        var leitores = Enumerable.Range(1, 3)
            .Select(id => LerAsync(id, canal.Reader, cts.Token))
            .ToArray();

        await Task.WhenAll(leitores);
        Console.WriteLine($"  tarefas ainda esperando: {leitores.Count(t => !t.IsCompleted)} de {leitores.Length}");
        relogio.Fim("  tempo");
    }

    private static async Task LerAsync(int id, ChannelReader<int> leitor, CancellationToken ct)
    {
        try
        {
            await leitor.ReadAsync(ct);
        }
        catch (OperationCanceledException)
        {
            Console.WriteLine($"  leitor {id}: saindo (OperationCanceledException)");
        }
    }
}
