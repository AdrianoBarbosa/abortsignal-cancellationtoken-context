using Cancelamento.Demos;

// Each example lives in Demos/ and runs on its own:
//
//   dotnet run -- 01
//
// Start the support server before the network examples (01 and 02):
//   node mock-server/server.js

var exemplos = new Dictionary<string, (string Descricao, Func<Task> Rodar)>
{
    ["01"] = ("Chamada HTTP cancelada por timeout", TimeoutHttp.RodarAsync),
    ["02"] = ("Cancelamento que some no meio do caminho (CancellationToken.None)", Propagacao.RodarAsync),
    ["03"] = ("Como o codigo e avisado: ThrowIfCancellationRequested e Register", Avisos.RodarAsync),
    ["04"] = ("Cancelamento em cascata com CreateLinkedTokenSource", Cascata.RodarAsync),
    ["05"] = ("Espera sem saida: ReadAsync de Channel sem token", EsperaSemSaida.RodarAsync),
};

var escolhido = args.FirstOrDefault();

if (escolhido is null || !exemplos.TryGetValue(escolhido, out var exemplo))
{
    Console.WriteLine("Exemplos disponiveis:\n");
    foreach (var (chave, (descricao, _)) in exemplos)
    {
        Console.WriteLine($"  dotnet run -- {chave}   {descricao}");
    }
    return;
}

Console.WriteLine($"== {escolhido}: {exemplo.Descricao} ==\n");
await exemplo.Rodar();
