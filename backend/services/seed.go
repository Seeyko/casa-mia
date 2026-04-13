package services

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

func SeedIfEmpty(db *Database, uploadDir string) {
	var count int
	err := db.DB.QueryRow(`SELECT COUNT(*) FROM locations`).Scan(&count)
	if err != nil || count > 0 {
		return
	}

	log.Println("Seeding database with initial data...")

	// === ADMIN USER ===
	adminUser := os.Getenv("ADMIN_USERNAME")
	adminPass := os.Getenv("ADMIN_PASSWORD")
	if adminUser == "" {
		adminUser = "admin"
	}
	if adminPass == "" {
		adminPass = "casamia2026"
		log.Println("WARNING: Using default admin password — set ADMIN_PASSWORD env var in production!")
	}
	hash, err2 := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
	if err2 == nil {
		db.DB.Exec(`INSERT INTO admin_users (username, password_hash) VALUES ($1, $2) ON CONFLICT (username) DO NOTHING`, adminUser, string(hash))
		log.Printf("Admin user created: %s", adminUser)
	}

	// === LOCATIONS ===
	db.DB.Exec(`INSERT INTO locations (name, slug, address, phone, opening_hours, order_method, order_info) VALUES
		('Entraigues-sur-la-Sorgue', 'entraigues', '51 Rue Laurent Bertrand, 84320 Entraigues-sur-la-Sorgue', '06 45 79 49 30',
		 '{"lundi":null,"mardi":{"slots":[{"open":"9h","close":"13h"},{"open":"16h","close":"21h"}]},"mercredi":{"slots":[{"open":"9h","close":"13h"},{"open":"16h","close":"21h"}]},"jeudi":{"slots":[{"open":"9h","close":"13h"},{"open":"16h","close":"21h"}]},"vendredi":{"slots":[{"open":"9h","close":"13h"},{"open":"16h","close":"21h30"}]},"samedi":{"slots":[{"open":"9h","close":"13h"},{"open":"16h","close":"21h30"}]},"dimanche":null}',
		 'sms', 'SMS ou WhatsApp au 06 45 79 49 30'),
		('Althen-des-Paluds', 'althen', '254 Avenue Ernest Perrin, 84210 Althen-des-Paluds', '04 90 36 16 33',
		 '{"lundi":{"slots":[{"open":"18h","close":"21h30"}]},"mardi":{"slots":[{"open":"18h","close":"21h30"}]},"mercredi":null,"jeudi":{"slots":[{"open":"18h","close":"21h30"}]},"vendredi":{"slots":[{"open":"18h","close":"22h"}]},"samedi":{"slots":[{"open":"18h","close":"22h"}]},"dimanche":{"slots":[{"open":"18h","close":"21h30"}]}}',
		 'phone', 'Appel uniquement au 04 90 36 16 33')
	`)

	// === MENU CATEGORIES (carte) ===
	db.DB.Exec(`INSERT INTO menu_categories (name, section, sort_order) VALUES
		('Pizzas Tomate', 'carte', 1),
		('Pizzas Creme', 'carte', 2),
		('Supplements', 'carte', 3),
		('Snacking', 'carte', 4),
		('Planches', 'carte', 5),
		('Pizza Dessert / Calzone', 'carte', 6),
		('Dolci', 'carte', 7),
		('Boissons', 'carte', 8)
	`)

	// === MENU CATEGORIES (traiteur) ===
	db.DB.Exec(`INSERT INTO menu_categories (name, section, sort_order) VALUES
		('Formules Aperitivo Italiano', 'traiteur', 1),
		('Coeurs d''Apero au Choix', 'traiteur', 2),
		('Bouchees', 'traiteur', 3),
		('Brochettes', 'traiteur', 4),
		('Verrines', 'traiteur', 5),
		('Desserts Traiteur', 'traiteur', 6)
	`)

	// === PIZZAS TOMATE (cat 1) ===
	db.DB.Exec(`INSERT INTO menu_items (category_id, name, description, price, sort_order, image_path) VALUES
		(1, 'Margherita', 'Tomate, parmesan, mozza, pesto', '9,50€', 1, ''),
		(1, 'Massilia', 'Tomate, moitie mozza parmesan / et moitie anchois', '10,00€', 2, ''),
		(1, 'Napolitaine', 'Tomate, parmesan, mozza, anchois, origan, capres', '10,50€', 3, ''),
		(1, 'Marinara', 'Tomate, parmesan, mozza, oignons, ail, tomates cerises', '10,50€', 4, ''),
		(1, 'Funghi', 'Tomate, parmesan, mozza, champignons, pesto', '10,50€', 5, ''),
		(1, 'Porchetta', 'Tomate, parmesan, mozza, porchetta (jambon blanc italien)', '10,50€', 6, ''),
		(1, 'Diavola', 'Tomate, parmesan, mozza, spianata (saucisson piquant italien)', '10,50€', 7, 'pizza-diavola-new.jpg'),
		(1, 'Regina', 'Tomate, parmesan, mozza, champignons, porchetta', '11,00€', 8, ''),
		(1, 'Formaggi', 'Tomate, parmesan, mozza, gorgonzola, taleggio, chevre', '12,00€', 9, ''),
		(1, 'Verdura', 'Tomate, parmesan, mozza, artichauts, poivrons, champignons, oignons rouges', '12,00€', 10, ''),
		(1, 'Salsiccia e Cipolle', 'Tomate, parmesan, mozza, saucisse sicilienne, oignons rouges', '13,00€', 11, ''),
		(1, 'Primavera', 'Tomate, parmesan, mozza, courgette, gorgonzola, speck', '13,00€', 12, ''),
		(1, '4 Stagioni', 'Tomate, parmesan, mozza, anchois, porchetta, artichauts, champignons', '13,00€', 13, 'pizza-quattro-stagioni.jpg'),
		(1, 'Calzone Salee', 'Tomate, parmesan, mozza, champignons, porchetta, oeuf', '13,00€', 14, ''),
		(1, 'San Daniele', 'Tomate, parmesan, mozza, San Daniele (jambon cru), roquette, tomates cerises', '13,00€', 15, ''),
		(1, 'Melenzana', 'Tomate, parmesan, mozza, aubergines, viande hachee de boeuf, tomates cerises', '13,00€', 16, ''),
		(1, 'Tuto Carni', 'Tomate, mozza, parmesan, viande hachee de boeuf, saucisse sicilienne, guanciale, spianata', '14,00€', 17, ''),
		(1, 'Burrata', 'Tomate, parmesan, mozza, roquette, tomates cerise, speck, burrata', '14,00€', 18, 'pizza-burrata.jpg'),
		(1, 'Tonno', 'Tomate, mozza, parmesan, pesto, thon, oignons rouge, tomate cerise', '14,00€', 19, 'pizza-tonno.jpg'),
		(1, 'Armenienne', 'Tomate, parmesan, origan, viande hachee, poivrons grilles, oignons rouges frais', '14,00€', 20, ''),
		(1, 'Marocchino', 'Tomate, parmesan, mozza, merguez, poivrons grilles, oignons rouges frais', '14,00€', 21, ''),
		(1, 'Corsica', 'Tomate, parmesan, mozza, origan, figatelli, ail confit', '14,00€', 22, ''),
		(1, 'Coppa Buffala', 'Tomate, parmesan, mozza, mozza di buffala, roquette, coppa, vinaigre balsamique', '15,00€', 23, '')
	`)

	// Update badges
	db.DB.Exec(`UPDATE menu_items SET badge = 'NEW' WHERE name IN ('Massilia', 'Primavera', 'Burrata') AND category_id = 1`)

	// === PIZZAS CREME (cat 2) ===
	db.DB.Exec(`INSERT INTO menu_items (category_id, name, description, price, sort_order, badge, note, image_path) VALUES
		(2, 'Capre e Miele', 'Creme, parmesan, mozza, chevre, miel, guanciale (lard italien), roquette', '12,00€', 1, '', '', ''),
		(2, 'Taleggio', 'Creme, parmesan, mozza, guanciale (lard italien), PDT sautees, oignons rouges, taleggio', '12,00€', 2, '', '', 'pizza-boscaiola-bacon-oignons.jpg'),
		(2, 'Bologne', 'Creme, pesto de pistaches, parmesan, mozza, mortadelle, eclats de pistaches', '13,00€', 3, '', '', ''),
		(2, 'Manzo', 'Creme, parmesan, mozza, gorgonzola, bresaola (boeuf seche), roquette', '13,00€', 4, '', '', 'pizza-bresaola.jpg'),
		(2, 'Friarielli', 'Creme, parmesan, mozza, friarielli (brocoli napolitain), saucisse sicilienne', '13,00€', 5, 'NEW', '', ''),
		(2, 'Pollo', 'Creme, mozza, parmesan, poulet, tomates cerises, pesto', '13,00€', 6, '', '', ''),
		(2, 'Raviole', 'Creme, mozza, parmesan, ravioles de Romans, jambon cru, pesto', '14,00€', 7, '', '', 'pizza-creme-jambon-cru-pesto.jpg'),
		(2, 'Raclette', 'Creme, mozza, parmesan, porchetta, fromage raclette, jambon cru (San Daniele)', '14,00€', 8, 'NEW', '', ''),
		(2, 'Salmone', 'Creme, parmesan, mozza, roquette, oignons rouge, saumon fume, vinaigre balsamique', '15,00€', 9, '', '', ''),
		(2, 'Porcini', 'Creme de cepes, mozza, parmesan, champignons, pancetta', '15,00€', 10, '', '', ''),
		(2, 'Carbonara', 'Mozza, guanciale (lard italien), parmesan, sauce carbo au jaune d''oeuf, parmesan et poivre noir', '15,00€', 11, '', '', ''),
		(2, 'Camembert au Four', 'Camembert cuit dans notre pate a pizza accompagne de charcuterie italienne', '16,00€', 12, 'NEW', '', 'planche-camembert-charcuterie.jpg'),
		(2, 'Tartufo', 'Creme de truffe, mozza, parmesan, magret fume, camembert di bufala', '18,00€', 13, '★', 'De mai a aout, possibilite de carpaccio de truffes fraiches d''ete — 5,00€', 'pizza-tartufo.jpg')
	`)

	// === SUPPLEMENTS (cat 3) ===
	db.DB.Exec(`INSERT INTO menu_items (category_id, name, description, price, sort_order) VALUES
		(3, 'Legume', '', '0,50€', 1),
		(3, 'Viande, charcuterie ou fromage', '', '1,50€', 2),
		(3, 'Burrata', '', '2,50€', 3),
		(3, 'Carpaccio de truffes d''ete fraiches', 'De mai a aout', '5,00€', 4)
	`)

	// === SNACKING (cat 4) ===
	db.DB.Exec(`INSERT INTO menu_items (category_id, name, description, price, sort_order) VALUES
		(4, 'Arancini Bolognaise ou Porchetta', '', '4,00€', 1),
		(4, 'Croq''Truffe', '', '5,00€', 2),
		(4, 'Piadine', 'Mozza, pesto, jambon cru', '3,50€', 3),
		(4, 'Panzo Napolitain', '', '5,00€', 4),
		(4, 'Focaccia Garnie', '', '4,00€', 5),
		(4, 'Bruschetta Sicilienne', '', '4,00€', 6)
	`)

	// === PLANCHES (cat 5) ===
	db.DB.Exec(`INSERT INTO menu_items (category_id, name, description, price, sort_order, image_path) VALUES
		(5, 'Charcuteries ou Charcuteries Fromages', '', '7,00€ / personne', 1, 'planche-charcuterie-premium.jpg'),
		(5, 'Planche Raclette', '200g de Taleggio DOP et assortiment de charcuterie italienne', '9,50€ / personne', 2, 'planche-charcuterie-fromages.jpg')
	`)

	// === PIZZA DESSERT (cat 6) ===
	db.DB.Exec(`INSERT INTO menu_items (category_id, name, description, price, sort_order) VALUES
		(6, 'Pizza Bueno', 'Duo de pate a tartiner choco-noisettes, Kit-Kat ball', '10,00€', 1),
		(6, 'Pizza Pistacchio', 'Duo de pate a tartiner choco-pistache, eclats de pistaches', '10,00€', 2),
		(6, 'Supplement Banane', '', '1,00€', 3)
	`)

	// === DOLCI (cat 7) ===
	db.DB.Exec(`INSERT INTO menu_items (category_id, name, description, price, sort_order, image_path) VALUES
		(7, 'Tiramisu traditionnel au cafe', '', '3,90€', 1, 'dolci-assortiment.jpg'),
		(7, 'Tiramisu chocolat noisettes et speculoos', '', '3,90€', 2, ''),
		(7, 'Panna cotta coulis de fruits rouges', '', '3,90€', 3, ''),
		(7, 'Panna cotta exotique', '', '3,90€', 4, '')
	`)

	// === BOISSONS (cat 8) ===
	db.DB.Exec(`INSERT INTO menu_items (category_id, name, description, price, sort_order) VALUES
		(8, 'Coca-Cola (1,5L)', '', '3,50€', 1),
		(8, 'Oasis (2L)', '', '3,50€', 2),
		(8, 'Estathe (the peche italien)', '', '3,00€', 3),
		(8, 'Biere Moretti Originale ou Blanche', '', '2,50€', 4),
		(8, 'Biere Moretti Rossa', '', '2,80€', 5),
		(8, 'Vin Sicilien Luna — Rouge', '', '13,00€', 6),
		(8, 'Vin Sicilien Luna — Blanc ou Rose', '', '11,00€', 7)
	`)

	// === TRAITEUR: FORMULES (cat 9) ===
	db.DB.Exec(`INSERT INTO menu_items (category_id, name, description, price, sort_order, note) VALUES
		(9, 'Primo', '1 Coeur d''apero au choix + Focaccia maison + 6 pieces au choix', '17€ / personne', 1, 'Minimum 6 personnes. Sur commande, minimum 3 jours a l''avance.'),
		(9, 'Secondo', '1 Coeur d''apero au choix + Focaccia maison + 9 pieces au choix', '22€ / personne', 2, 'Minimum 6 personnes. Sur commande, minimum 3 jours a l''avance.'),
		(9, 'Terso', '2 Coeurs d''apero au choix + Focaccia maison + 9 pieces au choix', '27€ / personne', 3, 'Minimum 6 personnes. Sur commande, minimum 3 jours a l''avance.')
	`)

	// === TRAITEUR: COEURS D'APERO (cat 10) ===
	db.DB.Exec(`INSERT INTO menu_items (category_id, name, description, price, sort_order, image_path) VALUES
		(10, 'Planche charcuteries et fromages italiens', '', '7,00€ / pers.', 1, ''),
		(10, 'Antipasti (assortiment de legumes grilles)', '', '6,50€ / pers.', 2, ''),
		(10, 'Assortiment de legumes croquants et trilogie de sauces maison', '', '5,50€ / pers.', 3, ''),
		(10, 'Saumon Gravelax', '', '7,00€ / pers.', 4, 'coeur-apero-saumon-gravelax.jpg'),
		(10, 'Carpaccio de Bresaola', '', '7,00€ / pers.', 5, 'coeur-apero-carpaccio-tomates.jpg')
	`)

	// === TRAITEUR: BOUCHEES (cat 11) ===
	db.DB.Exec(`INSERT INTO menu_items (category_id, name, description, price, sort_order, image_path) VALUES
		(11, 'Wrap facon Cesar', '', '2€ / piece', 1, ''),
		(11, 'Wrap thon', '', '2€ / piece', 2, ''),
		(11, 'Mini burger polpettes de poulet', '', '2€ / piece', 3, 'bouchees-mini-burgers-focaccia.jpg'),
		(11, 'Tramezini saumon fume', '', '2€ / piece', 4, ''),
		(11, 'Mini croq''Truffe', '', '2€ / piece', 5, ''),
		(11, 'Tourte Napolitaine', '', '2€ / piece', 6, 'bouchees-feuilletes-assortis.jpg'),
		(11, 'Quiche Italienne', 'Basilic, mozza, tomates confites', '2€ / piece', 7, ''),
		(11, 'Cake pesto mozza speck', '', '2€ / piece', 8, ''),
		(11, 'Cake tomates aubergines parmesan', '', '2€ / piece', 9, ''),
		(11, 'Navette saumon', '', '2€ / piece', 10, 'bouchees-navettes-saumon.jpg'),
		(11, 'Navette gorgonzola bresaola', '', '2€ / piece', 11, ''),
		(11, 'Bouchee cepes', '', '2€ / piece', 12, ''),
		(11, 'Bouchee truffe', '', '2€ / piece', 13, ''),
		(11, 'Bouchee jambon cru parmesan', '', '2€ / piece', 14, ''),
		(11, 'Croustillant mozza jambon cru', '', '2€ / piece', 15, ''),
		(11, 'Burger boeuf/porc', '', '2€ / piece', 16, ''),
		(11, 'Hot dog italien', '', '2€ / piece', 17, '')
	`)

	// === TRAITEUR: BROCHETTES (cat 12) ===
	db.DB.Exec(`INSERT INTO menu_items (category_id, name, description, price, sort_order, image_path) VALUES
		(12, 'Poulet, tomates confites', '', '', 1, ''),
		(12, 'Tomates cerises, mozza, speck', '', '', 2, 'brochettes-tomate-mozza-speck.jpg'),
		(12, 'Melon, jambon cru', '', '', 3, 'bouchees-melon-jambon.jpg'),
		(12, 'Bresaola, roquette, parmesan, tomates', '', '', 4, 'brochettes-bresaola-balsamique.jpg'),
		(12, 'Crevettes marinees basilic citron', '', '', 5, 'brochettes-crevettes.jpg'),
		(12, 'Legumes grilles', '', '', 6, ''),
		(12, 'Brochette magret figues', '', '', 7, '')
	`)

	// === TRAITEUR: VERRINES (cat 13) ===
	db.DB.Exec(`INSERT INTO menu_items (category_id, name, description, price, sort_order, image_path) VALUES
		(13, 'Pannacotta, courgette, chevre', '', '', 1, ''),
		(13, 'Tiramisu sale', '', '', 2, ''),
		(13, 'Mousse de saumon fume', '', '', 3, 'verrine-mousse-saumon.jpg'),
		(13, 'Brouillade a la truffe d''ete', '', '', 4, ''),
		(13, 'Soupe froide de tomates, basilic, coppa', '', '', 5, ''),
		(13, 'Tartare de boeuf a l''italienne', '', '', 6, 'verrine-tartare-boeuf.jpg'),
		(13, 'Tartare de crevettes a l''italienne', '', '', 7, ''),
		(13, 'Salade de poulpe', '', '', 8, ''),
		(13, 'Verrine facon Cesar', '', '', 9, '')
	`)

	// === TRAITEUR: DESSERTS (cat 14) ===
	db.DB.Exec(`INSERT INTO menu_items (category_id, name, description, price, sort_order, image_path) VALUES
		(14, 'Tiramisu cafe ou choco', '', '', 1, ''),
		(14, 'Pannacotta coulis au choix', 'Fruits rouges, pistaches, fruits exotiques', '', 2, ''),
		(14, 'Petit canolo', '', '', 3, 'dolci-mini-cannoli.jpg'),
		(14, 'Mini fondant chocolat', '', '', 4, ''),
		(14, 'Verrine baba au rhum ou limoncello', '', '', 5, '')
	`)

	// === COPY SEED IMAGES TO UPLOADS ===
	seedImages := []string{
		"pizza-diavola-new.jpg",
		"pizza-quattro-stagioni.jpg",
		"pizza-burrata.jpg",
		"pizza-tonno.jpg",
		"pizza-boscaiola-bacon-oignons.jpg",
		"pizza-bresaola.jpg",
		"pizza-creme-jambon-cru-pesto.jpg",
		"planche-camembert-charcuterie.jpg",
		"pizza-tartufo.jpg",
		"planche-charcuterie-premium.jpg",
		"planche-charcuterie-fromages.jpg",
		"dolci-assortiment.jpg",
		"coeur-apero-saumon-gravelax.jpg",
		"coeur-apero-carpaccio-tomates.jpg",
		"bouchees-mini-burgers-focaccia.jpg",
		"bouchees-feuilletes-assortis.jpg",
		"bouchees-navettes-saumon.jpg",
		"brochettes-tomate-mozza-speck.jpg",
		"bouchees-melon-jambon.jpg",
		"brochettes-bresaola-balsamique.jpg",
		"brochettes-crevettes.jpg",
		"verrine-mousse-saumon.jpg",
		"verrine-tartare-boeuf.jpg",
		"dolci-mini-cannoli.jpg",
	}
	copySeedImages(seedImages, uploadDir)

	log.Println("Database seeded successfully")
}

// copySeedImages copies seed photos from frontend/images/ to the backend uploads directory.
// It looks for the frontend images relative to the working directory or common project layouts.
func copySeedImages(images []string, uploadDir string) {
	// Try multiple possible paths for the frontend images directory
	candidates := []string{
		filepath.Join("/app", "seed-images"),
		filepath.Join("..", "frontend", "images"),
		filepath.Join("frontend", "images"),
		filepath.Join("..", "..", "frontend", "images"),
	}

	var srcDir string
	for _, c := range candidates {
		abs, err := filepath.Abs(c)
		if err != nil {
			continue
		}
		if info, err := os.Stat(abs); err == nil && info.IsDir() {
			srcDir = abs
			break
		}
	}

	if srcDir == "" {
		log.Println("WARNING: Could not find frontend/images/ directory for seed images")
		return
	}

	// Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("WARNING: Could not create upload directory %s: %v", uploadDir, err)
		return
	}

	absUploadDir, _ := filepath.Abs(uploadDir)
	copied := 0

	for _, img := range images {
		srcPath := filepath.Join(srcDir, img)
		dstPath := filepath.Join(absUploadDir, img)

		// Skip if destination already exists
		if _, err := os.Stat(dstPath); err == nil {
			continue
		}

		if err := copyFile(srcPath, dstPath); err != nil {
			log.Printf("WARNING: Could not copy seed image %s: %v", img, err)
		} else {
			copied++
		}
	}

	if copied > 0 {
		log.Printf("Copied %d seed images to %s", copied, absUploadDir)
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
