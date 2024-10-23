import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import { SidebarDemo } from "./ui/sidebar";
import { motion } from "framer-motion";
import {
  IconHeart,
  IconMessageCircle,
  IconCalendarEvent,
  IconPlus,
  IconBell,  // Ajout de l'icône Bell
  IconMessage,  // Ajout de l'icône Chat
  IconUsersGroup // Ajout de l'icône Groupes
} from "@tabler/icons-react";
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
    <div className="flex flex-col lg:flex-row min-h-screen bg-gray-900 relative">
      <SidebarDemo />

      {/* Main Content */}
      <div className="flex-1 p-4 lg:p-8 bg-gray-800 shadow-lg rounded-lg m-5 relative">

        {/* Global Header with Notifications, Chat, and Groups (Fixed Position) */}
        <div className="fixed top-4 right-4 flex space-x-4 z-50 bg-gray-800 p-2 rounded-lg"> {/* Ajout de z-index et de fond */}
          <button className="text-white">
            <IconBell className="h-6 w-6" />
          </button>
          <button className="text-white">
            <IconMessage className="h-6 w-6" />
          </button>
          <button className="text-white">
            <IconUsersGroup className="h-6 w-6" />
          </button>
        </div>

        {/* Quick Action Cards */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8 mt-12"> {/* Ajustement de l'espacement avec mt-12 */}
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

          {/* Créer un post */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            className="p-6 bg-green-600 text-white rounded-lg shadow-md hover:bg-green-500 transition duration-300 flex items-center"
          >
            <IconPlus className="h-8 w-8 mr-3" />
            <div>
              <h3 className="text-xl font-bold mb-2">Créer un Post</h3>
              <p>Partagez vos pensées avec le réseau.</p>
            </div>
          </motion.div>

          {/* Créer un événement */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.7 }}
            className="p-6 bg-purple-600 text-white rounded-lg shadow-md hover:bg-purple-500 transition duration-300 flex items-center"
          >
            <IconPlus className="h-8 w-8 mr-3" />
            <div>
              <h3 className="text-xl font-bold mb-2">Créer un Événement</h3>
              <p>Organisez un événement pour votre groupe ou réseau.</p>
            </div>
          </motion.div>
        </div>

        {/* Recent Posts */}
        <div className="bg-gray-700 p-4 lg:p-6 rounded-lg shadow-md">
          <h2 className="text-2xl font-bold text-gray-300 mb-4">Vos Posts Récents</h2>
          <div className="space-y-4 lg:space-y-6">
            {[1, 2, 3].map((post, idx) => (
              <div key={idx} className="p-6 bg-gray-800 shadow-lg rounded-lg hover:shadow-xl transition duration-300">
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center space-x-4">
                    <Image src="/avatar.jpg" alt="User avatar" width={50} height={50} className="rounded-full" />
                    <div>
                      <h3 className="text-lg font-semibold text-gray-300">Username</h3>
                      <p className="text-gray-400 text-sm">Posté le 23 octobre 2024, 15h45</p>
                    </div>
                  </div>
                </div>

                <div className="mb-4">
                  <p className="text-gray-400">Ceci est le contenu du post {post}. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque at est non erat commodo facilisis.</p>
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
