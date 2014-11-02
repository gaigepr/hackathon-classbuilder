import java.util.*;
import java.io.*;

public class ClassParser {
	static File parseFile;
	static String nameString;
	static int reqCount;

	public static void main(String[] args) throws FileNotFoundException, IOException{
		if( args.length != 1) {
			System.err.println("Error: Requires one args");
			System.exit(1);
		}
      
        String[] nameParsed = args[0].split("\\.");
        if (true) {
        }

		PrintWriter writer = new PrintWriter(nameParsed[0], "UTF-8");


		parseFile = new File(args[0]);	
		Scanner input = new Scanner(parseFile);
		Scanner titleScan = new Scanner(parseFile);
		input.useDelimiter("<h3>");
		titleScan.useDelimiter("<h3>");

		String department;
		String course;
		String courseNum;
		String credits;
		String title;
		String description;
		String gurAttr;
		String prereq;
		String name;
		String num;
		ArrayList<String> namePool = new ArrayList<String>();

		input.next();
		titleScan.next();

		while(titleScan.hasNext()) {
			String classChunk = titleScan.next();

			course = getCourse(classChunk);
			String[] nameSplit = course.split("\\s+");
			name = nameSplit[0];
			namePool.add(name);
		}
		namePool.add("ETEC");
		namePool.add("SCED");
		namePool.remove("A");

		nameString = namePool.toString();

		while(input.hasNext()) {
			String classChunk = input.next();

			course = getCourse(classChunk);
			String[] courseSplit = course.split("\\s+");
			String[] titleSplit = course.split("\\s+-\\s+");
			department = courseSplit[0];
			courseNum = courseSplit[1];
         if (titleSplit.length == 1) {
            title = "***TITLE ERROR***";
         } else {
			   title = titleSplit[1];
         }

			credits = getCredits(classChunk);
			description = getDescription(classChunk);
			gurAttr = getGUR(classChunk);
			prereq = getPrereq(classChunk);

			writer.println(department);
			writer.println(courseNum);
			writer.println(credits);
			writer.println(title);
			writer.println(description);
			writer.println(gurAttr);
			writer.println(prereq);
         writer.println("");
         writer.println("");
         writer.println("");
		}

		writer.close();
	}

	public static String getCourse(String cChunk) {
		String hClose = "</h3>";
		String lineDep = "";

		Scanner lineData = new Scanner(cChunk);
		lineData.useDelimiter(hClose);
		lineDep = lineData.next();

		return lineDep;
	}

	public static String getCredits(String cChunk) {
		String lineDep = "";

		Scanner lineData = new Scanner(cChunk);
		lineData.useDelimiter("Credits: ");
		lineData.next();
		lineDep = lineData.next();

		return lineDep.substring(0, 1);
	}

	public static String getDescription(String cChunk) {
		String lineDep = "";

		Scanner lineData = new Scanner(cChunk);
		lineData.useDelimiter("<hr>");
		lineData.next();
		lineDep = lineData.next();

		String descrSplit[] = lineDep.split("<br>");
		return descrSplit[0].replace("\n", "").replace("\r", "");
	}

		public static String getGUR(String cChunk) {
		String lineDep = "";

		if (cChunk.contains("GUR Attributes: ")) {
			Scanner lineData = new Scanner(cChunk);
			lineData.useDelimiter("GUR Attributes: ");
			lineData.next();
			lineDep = lineData.next();

			int i = 1;
			while (!lineDep.substring(0 + i, 1 + i).equals("<")) {
				i++;
			}

			return lineDep.substring(0, i);
		} else {
			return "0";
		}
	}

	public static String getPrereq(String cChunk) {
		String lineDep = "";
		int counter = 0;
		int i = 0;
		int pointer = 0;
		reqCount = 0;


		if(cChunk.contains("Prerequisites &amp; Notes:</strong><br>")) {
			Scanner lineData = new Scanner(cChunk);
			lineData.useDelimiter("Prerequisites &amp; Notes:</strong><br>");
			lineData.next();
			lineDep = lineData.next();

			String descrSplit[] = lineDep.split("<br>");
			String parsedString = descrSplit[0].replace("\n", "").replace("\r", "").replace(",", "")
			.replace(".", "").replace(":", "").replace("AND", "").replace(";", "").replace("Any of", "ONE")
			.replace("any of", "ONE").replace("One of", "ONE").replace("one of", "ONE").replace("or equivalent", "")
			.replace("or concurrent", "").replace("or higher", "").replace("or", "OR").trim();

			if (parsedString.contains("Department") || parsedString.contains("Permission") || parsedString.contains("permission")) {
				return "1[CSCI 000]";
			} else if (parsedString.contains("Final")) {
				return "1[CSCI 999]";
			} else if (parsedString.contains("Graduate")) {
				return "1[CSCI 888]";
			} else {
				String[] parsedStringSplit = parsedString.split("\\s+");
				int parsedLen = parsedStringSplit.length;
				String retString = "";
				String tempString = "";

				if (!parsedString.contains("ONE") && !parsedString.contains("OR")) {
					while(i < parsedLen) {
						if (nameString.contains(parsedStringSplit[i])) {
							counter++;
						}
						i++;					
					}

					String[] strArray = new String[counter];

					i = 0;
					counter = 0;
					while(i < parsedLen) {
						if (nameString.contains(parsedStringSplit[i])) {
							strArray[counter] = parsedStringSplit[i] + " " + parsedStringSplit[i + 1];
							i++;
							counter++;
						}
						i++;
					}

					retString += counter + "[";
					int len = strArray.length;	
					for (i = 0; i < len - 1; i++) {
						retString += strArray[i] + ",";
					}

					retString += strArray[len - 1] + "]";

				} else {
					ArrayList<String> tempAL = new ArrayList<String>();

					i = 0;
					while(i < parsedLen) {
						if (nameString.contains(parsedStringSplit[i])) {
							tempAL.add(parsedStringSplit[i] + " " + parsedStringSplit[i + 1]);
							i++;
						}

						if (parsedStringSplit[i].equals("OR")) {
							tempAL.add("OR");

						}

						if (parsedStringSplit[i].equals("ONE")) {
							tempAL.add("ONE");

						}
						i++;
					}
					reqCount = tempAL.size();

					while(tempAL.contains("ONE")) {
						tempAL = mergeClassesONE(tempAL);
					}

					while(tempAL.contains("OR")) {
						tempAL = mergeClassesOR(tempAL);
					}

					retString = tempAL.toString();

					retString = reqCount + retString;

				}
				return retString;
			}
		} else {
			return "0[]";
		}
	}

	public static ArrayList<String> mergeClassesOR(ArrayList<String> listArray) {
		int pointer = listArray.indexOf("OR");


		while(pointer == 0) {
			listArray.remove(pointer);
			pointer = listArray.indexOf("OR");
			reqCount--;
		}

		int size = listArray.size();
      if (listArray.contains("OR")) {
   		while (listArray.get(size - 1).equals("OR")) {
   			listArray.remove(size - 1);
   			size = listArray.size();
   			reqCount--;
   		}
      }

		if (listArray.contains("OR")) {
			String combine = "1[" + listArray.get(pointer - 1) + "," + listArray.get(pointer + 1) + "]";
			pointer--;
			listArray.remove(pointer);
			listArray.remove(pointer);
			listArray.remove(pointer);
			reqCount -= 2;

			listArray.add(pointer, combine);
		}

		return listArray;
	}

	public static ArrayList<String> mergeClassesONE(ArrayList<String> listArray) {
		int pointer = listArray.indexOf("ONE");
		reqCount = pointer + 1;
		int end = listArray.size();

		String combine = "1[";

		for (int i = pointer + 1; i < end - 1; i++) {
			combine += listArray.get(i) + ",";
		}

		combine += listArray.get(end - 1) + "]";

		while(end > pointer) {
			listArray.remove(pointer);
         end--;
		}

		listArray.add(pointer, combine);	
		return listArray;
	}
}