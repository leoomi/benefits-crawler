import React from 'react';
import Form from './components/form'
import Results from './components/results';

function App() {
  const [results, setResults] = React.useState<string[] | null>(null);
  const [submitted, setSubmitted] = React.useState<boolean>(false);

  const handleSubmit = (event: React.SyntheticEvent) => {
    event.preventDefault();
    const target = event.target as typeof event.target & {
      cpf: { value: string };
    };

    const cpf = target.cpf.value;
    setSubmitted(true);
    console.log(cpf);
  }

  return (
    <div className="fixed h-screen w-screen flex flex-col justify-start items-center bg-gray-900 text-white shadow-lg">
      <h1 className="text-6xl m-8">Benefícios</h1>
      <div className="w-11/12 h-96 bg-gray-800 rounded-lg shadow-lg">
        <Form onSave={handleSubmit} />
        <Results submitted={submitted} results={results}/>
      </div>
    </div>
  );
}

export default App;
