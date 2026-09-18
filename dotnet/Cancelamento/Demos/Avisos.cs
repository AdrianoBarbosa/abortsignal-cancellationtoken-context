namespace Cancelamento.Demos;

/// <summary>
/// Example 3: how running code finds out it should stop.
///
/// Two ways: checking the token between steps (polling), or registering a
/// callback with Register. The callback is how you stop something that knows
/// nothing about CancellationToken, like the legacy timer below.
/// </summary>
public static class Avisos
{
    public static async Task RodarAsync()
    {
        await ComPollingAsync();
        Console.WriteLine();
        await ComRegisterAsync();
    }

    private static async Task ComPollingAsync()
    {
        Console.WriteLine("-- polling: ThrowIfCancellationRequested entre um lote e outro --");
        var relogio = Apoio.Cronometrar();

        using var cts = new CancellationTokenSource(TimeSpan.FromMilliseconds(250));
        var lotes = 0;

        try
        {
            while (true)
            {
                cts.Token.ThrowIfCancellationRequested();
                await Task.Delay(50); // one batch of work, deliberately without the token
                lotes++;
            }
        }
        catch (OperationCanceledException)
        {
            Console.WriteLine($"  parou depois de {lotes} lotes: OperationCanceledException");
        }

        relogio.Fim("  tempo");
    }

    private static async Task ComRegisterAsync()
    {
        Console.WriteLine("-- callback: Register parando uma tarefa legada --");
        var relogio = Apoio.Cronometrar();

        using var cts = new CancellationTokenSource(TimeSpan.FromMilliseconds(350));

        // A legacy task: it only knows how to be disposed, it takes no token.
        var ticks = 0;
        var legado = new Timer(
            _ => Console.WriteLine($"  tarefa legada: tick {Interlocked.Increment(ref ticks)}"),
            null,
            dueTime: 100,
            period: 100);

        // Without RunContinuationsAsynchronously, the rest of this method would run
        // inline inside the Register callback, on the thread that canceled.
        var parou = new TaskCompletionSource(TaskCreationOptions.RunContinuationsAsynchronously);
        using var registro = cts.Token.Register(() =>
        {
            legado.Dispose();
            Console.WriteLine("  Register: parando a tarefa legada");
            parou.SetResult();
        });

        await parou.Task;
        relogio.Fim("  tempo");

        var ticksAoParar = Volatile.Read(ref ticks);
        await Task.Delay(300);
        Console.WriteLine($"  ticks depois de parar: {Volatile.Read(ref ticks) - ticksAoParar}");
    }
}
