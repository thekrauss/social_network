import React, { useState, useEffect } from "react";
import { useRouter } from "next/router";
import AOS from "aos";
import "aos/dist/aos.css";

export default function AuthPage() {
  const [isLogin, setIsLogin] = useState(true); // Gère le formulaire à afficher (connexion ou inscription)
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [formData, setFormData] = useState({
    username: "",
    age: "",
    email: "",
    password: "",
    firstName: "",
    lastName: "",
    gender: "",
    dateOfBirth: "",
    avatar: "",
    bio: "",
    phoneNumber: "",
    address: "",
    isPrivate: false,
  });
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const router = useRouter();

  useEffect(() => {
    AOS.init({ duration: 1000 });
  }, []);

  // Gère les changements pour Register
  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData({
      ...formData,
      [name]: type === "checkbox" ? checked : value,
    });
  };

  // Validation de l'email
  const validateEmail = (email) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  // Gère la soumission du formulaire Login
  const handleLoginSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");
  
    if (!validateEmail(email)) {
      setError("Veuillez entrer un email valide.");
      return;
    }
  
    try {
      const response = await fetch("http://localhost:8079/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }), // Envoi sous forme de JSON
      });
  
      if (!response.ok) {
        const errorResponse = await response.json();
        setError(`Erreur: ${errorResponse.message || "Échec de la connexion, vérifiez vos identifiants."}`);
      } else {
        const data = await response.json();
        localStorage.setItem("authToken", data.token); // Stockage du token JWT
        setSuccess("Connexion réussie !");
        setTimeout(() => {
          router.push("/"); // Redirection après succès
        }, 2000);
      }
    } catch (error) {
      console.log("Erreur réseau : ", error);
      setError("Erreur réseau, veuillez réessayer.");
    }
  };
  
  

  // Gère la soumission du formulaire Register
  const handleRegisterSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    // Validation de l'email et des autres champs à ajouter ici

    try {
      const response = await fetch("http://localhost:8079/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        const errorResponse = await response.json();
        setError(`Erreur: ${errorResponse.message || "Une erreur est survenue"}`);
      } else {
        setSuccess("Inscription réussie !");
        setFormData({
          username: "",
          age: "",
          email: "",
          password: "",
          firstName: "",
          lastName: "", 
          gender: "",
          dateOfBirth: "",
          avatar: "",
          bio: "",
          phoneNumber: "",
          address: "",
          isPrivate: false,
        });
        setTimeout(() => {
          router.push("/login");
        }, 2000);
      }
    } catch (error) {
      setError("Erreur réseau, veuillez réessayer.");
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 dark:bg-gray-900 text-gray-100 p-4 sm:p-10">
      <div data-aos="fade-up" className="w-full max-w-lg bg-gray-800 rounded-lg shadow-xl p-6 sm:p-10">
        <div className="text-center mb-6">
          <button
            onClick={() => setIsLogin(!isLogin)}
            className="text-indigo-400 hover:underline focus:outline-none"
          >
            {isLogin ? "Vous n'avez pas de compte ? Créez un compte" : "Déjà un compte ? Connectez-vous"}
          </button>
        </div>

        {/* Formulaire Login */}
        {isLogin ? (
          <>
            <h2 className="text-3xl sm:text-4xl font-extrabold text-center mb-6 animate-pulse">Connexion</h2>

            {error && <p className="text-red-500 text-center">{error}</p>}
            {success && <p className="text-green-500 text-center">{success}</p>}

            <form className="space-y-4 sm:space-y-6" onSubmit={handleLoginSubmit}>
              <div className="relative group">
                <label htmlFor="email" className="block text-sm font-medium text-gray-400">
                  Adresse email
                </label>
                <input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="mt-1 block w-full px-4 py-2 border border-gray-600 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 transition duration-300 bg-gray-700 text-white"
                  placeholder="Entrez votre email"
                  required
                />
              </div>

              <div className="relative group">
                <label htmlFor="password" className="block text-sm font-medium text-gray-400">
                  Mot de passe
                </label>
                <input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="mt-1 block w-full px-4 py-2 border border-gray-600 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 transition duration-300 bg-gray-700 text-white"
                  placeholder="Entrez votre mot de passe"
                  required
                />
              </div>

              <button type="submit" className="w-full py-3 px-4 bg-indigo-600 text-white rounded-lg shadow-md hover:bg-indigo-700 transition duration-300">
                Connexion
              </button>
            </form>
          </>
        ) : (
          <>
            {/* Formulaire Register */}
            <h2 className="text-3xl sm:text-4xl font-extrabold text-center mb-6 animate-pulse">Inscription</h2>

            {error && <p className="text-red-500 text-center">{error}</p>}
            {success && <p className="text-green-500 text-center">{success}</p>}

            <form className="space-y-4 sm:space-y-6" onSubmit={handleRegisterSubmit}>
              {/* Username */}
              <div className="relative group">
                <label htmlFor="username" className="block text-sm font-medium text-gray-400">
                  Nom d'utilisateur
                </label>
                <input
                  id="username"
                  name="username"
                  type="text"
                  value={formData.username}
                  onChange={handleChange}
                  className="mt-1 block w-full px-4 py-2 border border-gray-600 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 transition duration-300 bg-gray-700 text-white"
                  placeholder="Entrez votre nom d'utilisateur"
                  required
                />
              </div>

              {/* Age */}
              <div className="relative group">
                <label htmlFor="age" className="block text-sm font-medium text-gray-400">
                  Âge
                </label>
                <input
                  id="age"
                  name="age"
                  type="number"
                  value={formData.age}
                  onChange={handleChange}
                  className="mt-1 block w-full px-4 py-2 border border-gray-600 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 transition duration-300 bg-gray-700 text-white"
                  placeholder="Entrez votre âge"
                  required
                />
              </div>

              {/* Email */}
              <div className="relative group">
                <label htmlFor="email" className="block text-sm font-medium text-gray-400">
                  Adresse email
                </label>
                <input
                  id="email"
                  name="email"
                  type="email"
                  value={formData.email}
                  onChange={handleChange}
                  className="mt-1 block w-full px-4 py-2 border border-gray-600 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 transition duration-300 bg-gray-700 text-white"
                  placeholder="Entrez votre email"
                  required
                />
              </div>

              {/* Password */}
              <div className="relative group">
                <label htmlFor="password" className="block text-sm font-medium text-gray-400">
                  Mot de passe
                </label>
                <input
                  id="password"
                  name="password"
                  type="password"
                  value={formData.password}
                  onChange={handleChange}
                  className="mt-1 block w-full px-4 py-2 border border-gray-600 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 transition duration-300 bg-gray-700 text-white"
                  placeholder="Entrez votre mot de passe"
                  required
                />
              </div>

              <button type="submit" className="w-full py-3 px-4 bg-indigo-600 text-white rounded-lg shadow-md hover:bg-indigo-700 transition duration-300">
                S'inscrire
              </button>
            </form>
          </>
        )}
      </div>
    </div>
  );
}
