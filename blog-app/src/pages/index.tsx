import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import { SidebarDemo } from "./ui/sidebar.js"; 
import { motion } from "framer-motion";

export default function HomePage() {
  const router = useRouter();
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem("authToken");

    //  existe, l'utilisateur est authentifié
    if (token) {
      setIsAuthenticated(true);
    } else {
      // vers la page de login si non authentifié
      router.push("/login");
    }
  }, []);

  if (!isAuthenticated) {
    return (
      <div className="flex justify-center items-center h-screen">
        {/* Animation en attendant la redirection */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.5 }}
          className="text-lg font-bold text-gray-500"
        >
          Redirection...
        </motion.div>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen">
      {/* Sidebar intégré */}
      <SidebarDemo />
      {/* Contenu principal */}
      <div className="flex-1 p-8 bg-gray-100">
        <h1 className="text-3xl font-bold mb-4">Bienvenue sur la page d'accueil</h1>
        <p>le contenu l' application.</p>
      </div>
    </div>
  );
}
