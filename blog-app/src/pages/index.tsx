import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import { SidebarDemo } from "./ui/sidebar";
import { motion } from "framer-motion";
import { IconHeart, IconMessageCircle, IconSearch, IconCalendarEvent, IconMessage, IconChartLine } from "@tabler/icons-react";
import Image from "next/image";

export default function HomePage() {
  const router = useRouter();
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem("authToken");
    if (token) {
      setIsAuthenticated(true);
    } else {
      router.push("/login");
    }
  }, [router]);

  if (!isAuthenticated) {
    return (
      <div className="flex justify-center items-center h-screen bg-gray-900">
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ duration: 0.5 }} className="text-lg font-bold text-gray-500">
          Redirection...
        </motion.div>
      </div>
    );
  }

  return (
    <div className="flex flex-col lg:flex-row min-h-screen bg-gray-900">
      <SidebarDemo />
      <div className="flex-1 p-4 lg:p-8 bg-gray-800 shadow-lg rounded-lg m-5">
        <div className="mb-6 flex">
          <input
            type="text"
            placeholder="Rechercher..."
            className="flex-1 p-4 text-gray-300 bg-gray-700 rounded-lg shadow focus:outline-none focus:ring-2 focus:ring-cyan-500"
          />
          <button className="ml-3 p-4 bg-cyan-600 rounded-lg text-white shadow hover:bg-cyan-500 transition duration-300">
            <IconSearch className="h-5 w-5" />
          </button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            className="p-6 bg-cyan-600 text-white rounded-lg shadow-md hover:bg-cyan-500 transition duration-300 flex items-center"
          >
            <IconCalendarEvent className="h-8 w-8 mr-3" />
            <div>
              <h3 className="text-xl font-bold mb-2">Événements</h3>
              <p>Découvrez les événements à venir dans votre réseau.</p>
            </div>
          </motion.div>
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            className="p-6 bg-blue-600 text-white rounded-lg shadow-md hover:bg-blue-500 transition duration-300 flex items-center"
          >
            <IconMessage className="h-8 w-8 mr-3" />
            <div>
              <h3 className="text-xl font-bold mb-2">Messages</h3>
              <p>Consultez vos conversations récentes et envoyez des messages.</p>
            </div>
          </motion.div>
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.7 }}
            className="p-6 bg-purple-600 text-white rounded-lg shadow-md hover:bg-purple-500 transition duration-300 flex items-center"
          >
            <IconChartLine className="h-8 w-8 mr-3" />
            <div>
              <h3 className="text-xl font-bold mb-2">Statistiques</h3>
              <p>Visualisez vos performances récentes.</p>
            </div>
          </motion.div>
        </div>

        <div className="bg-gray-700 p-4 lg:p-6 rounded-lg shadow-md">
          <h2 className="text-2xl font-bold text-gray-300 mb-4">Vos Posts Récents</h2>
          <div className="space-y-4 lg:space-y-6">
            {[1, 2, 3].map((post, idx) => (
              <div key={idx} className="p-6 bg-gray-800 shadow-lg rounded-lg hover:shadow-xl transition duration-300">
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center space-x-4">
                    <Image
                      src="/avatar.jpg"
                      alt="User avatar"
                      width={50}
                      height={50}
                      className="rounded-full"
                    />
                    <div>
                      <h3 className="text-lg font-semibold text-gray-300">Username</h3>
                      <p className="text-gray-400 text-sm">Posté le 23 octobre 2024, 15h45</p>
                    </div>
                  </div>
                </div>

                <div className="mb-4">
                  <p className="text-gray-400">
                    Ceci est le contenu du post {post}. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque at est non erat commodo facilisis.
                  </p>
                </div>

                <div className="flex items-center mt-4 space-x-4">
                  <button className="flex items-center text-gray-400 hover:text-cyan-400">
                    <IconHeart className="h-5 w-5 mr-1" /> J'aime <span className="ml-2">(23)</span>
                  </button>
                  <button className="flex items-center text-gray-400 hover:text-cyan-400">
                    <IconMessageCircle className="h-5 w-5 mr-1" /> Commenter <span className="ml-2">(12)</span>
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}