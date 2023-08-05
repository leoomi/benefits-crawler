interface ResultsInput {
  submitted: boolean;
  results: string[] | null;
}

function CPFNotRegistered() {
  return <span>CPF não está na base de dados.</span>
}

function NoBefits() {
  return <span>Não há benefícios associados a este CPF.</span>
}

function Benefits(results: string[]) {
  return(
    results.map((r, idx) => 
      <span key={idx}>{r}</span>
    )
  )
}

export default function Results({ submitted, results }: ResultsInput) {
  let resultOption
  if (results === null) {
    resultOption = CPFNotRegistered();
  } else if (results.length === 0) {
    resultOption = NoBefits();
  } else {
    resultOption = Benefits(results);
  }

  return (
    <>
    {
      submitted &&
      <div className="flex flex-col items-center mt-8 mx-4">
        <h2 className="text-xl">Resultados</h2>
        { resultOption }
      </div>
    }
    </>
  )
}