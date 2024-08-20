import React, { useState } from 'react';
import TextareaEntrada from './Componentes/TextAreaEntrada';
import Botones from './Componentes/Botones';
import TextareaSalida from './Componentes/TextAreaSalida';
import './App.css';

const App: React.FC = () => {
  const [comandos, setComandos] = useState<string>('');
  const [resultado, setResultado] = useState<string>(''); // Inicializa como una cadena vacía

  const handleArchivoCargado = (contenido: string) => {
    setComandos(contenido);
  };

  const ejecutarComandos = async () => {
    try {
      const response = await fetch('http://localhost:3000/api', {
        method: 'POST',
        headers: {
          'Content-Type': 'text/plain',
        },
        body: comandos
      });

      const data = await response.json();
      setResultado(data.resultado || ''); // Usa una cadena vacía si no hay resultado
    } catch (error) {
      setResultado(`Error: ${error}`);
    }
  };

  return (
    <div className="app-container">
      <h1>Proyecto1</h1>
      <TextareaEntrada value={comandos} onChange={setComandos} />
      <Botones onArchivoCargado={handleArchivoCargado} onEjecutar={ejecutarComandos} />
      <TextareaSalida value={resultado} />
    </div>
  );
};

export default App;
