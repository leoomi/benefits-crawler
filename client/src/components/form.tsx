import React from "react";

interface FormInput {
  onSave(event: React.SyntheticEvent): void;
}

export default function Form({ onSave }: FormInput) {
  return (<form onSubmit={onSave}>
    <div className="flex flex-col items-center">
      <label className="m-4">
        <span className="mr-4">CPF:</span>
        <input className="px-4 py-2 text-black rounded-md"
          name="cpf"
          type="text" />
      </label>
      <button
        className="bg-gray-700 rounded-lg px-4 py-2 shadow-lg"
        type="submit"
        value="Submit">
        Buscar
      </button>
    </div>
  </form>
  )
}