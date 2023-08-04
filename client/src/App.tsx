import React from 'react';
import logo from './logo.svg';
import './App.css';

function App() {

  const handleSubmit = (event: React.SyntheticEvent) => {
    event.preventDefault();
    const target = event.target as typeof event.target & {
      cpf: { value: string };
    };

    const cpf = target.cpf.value;
    console.log(cpf);
  }

  return (
    <div className="fixed h-screen w-screen flex flex-col justify-start items-center bg-gray-900 text-white shadow-lg">
      <h1 className="text-6xl m-8">Benefícios</h1>
      <div className="w-11/12 h-96 bg-gray-800 rounded-lg shadow-lg">
        <form onSubmit={handleSubmit}>
          <div className="flex flex-col items-center">
            <label className="m-4">
              <span className="mr-4">CPF:</span>
              <input className="text-black rounded-md" 
                name="cpf"
                type="text"/>
            </label>
            <button 
              className="bg-gray-700 rounded-lg p-4 shadow-lg"
              type="submit"
              value="Submit">
                Buscar
              </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default App;
